package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
	"github.com/royanqodri/Login-Gateway-API/model/request"
	_ "github.com/royanqodri/Login-Gateway-API/model/response"
	"github.com/royanqodri/Login-Gateway-API/service"
	"github.com/royanqodri/Login-Gateway-API/util"
)

type SessionController interface {
	GetBySessionData(ctx *gin.Context)
}

type SessionControllerImpl struct {
	sessionService service.SessionService
}

func NewSessionController(sessionService service.SessionService) SessionController {
	return &SessionControllerImpl{sessionService: sessionService}
}

// GetBySessionData godoc
// @Summary Get Session Data
// @Description Get Session Data
// @Tags Session
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} response.SessionDataMainResponse
// @Failure 400 {object} response.MainResponse
// @Failure 500 {object} response.MainResponse
// @Router /session-data [get]
func (controller *SessionControllerImpl) GetBySessionData(ctx *gin.Context) {
	req := request.SessionParamRequest{}

	// Validate request payload
	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	// Extract token from Authorization header
	tokenString := util.GetTokenFromHeader(ctx)
	if tokenString == "" {
		util.WriteResponse(ctx, util.JSON, http.StatusUnauthorized, []any{}, []any{}, "Authorization header is missing or invalid", 0, 0)
		ctx.Abort()
		return
	}

	// Parse and validate token
	claims, err := util.ParseToken(tokenString)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusUnauthorized, []any{}, []any{}, fmt.Sprintf("invalid token: %v", err), 0, 0)
		ctx.Abort()
		return
	}

	customerNo := claims["customer_no"].(string)
	customerId := int64(claims["customer_id"].(float64))
	username := claims["username"].(string)
	site := req.Site
	module := claims["module"].(string)

	dataCustomer := request.SessionRequest{
		Username:   username,
		CustomerID: customerId,
		CustomerNo: customerNo,
		Site:       site,
		Menu:       module,
	}

	// Get data from service
	resp, err := controller.sessionService.GetBySessionData(ctx, dataCustomer)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("get by session data failed: %v", err), 0, 0)
		return
	}

	// Build Redis request payload
	req2 := request.RedisRequest{
		Username:   username,
		CustomerID: customerId,
		CustomerNo: customerNo,
		Menu:       "",
	}

	// Fetch user modules
	userModules, err := controller.sessionService.GetUserModules(ctx, req2)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("Error fetching user permissions: %v", err), 0, 0)
		return
	}

	// Serialize the user modules (permissions) into JSON
	modulesJson, err := json.Marshal(userModules)
	if err != nil {
		log.Printf("Error marshaling user modules: %v", err)
	}

	// Store the token and the corresponding permission data in Redis
	expUnix := int64(claims["exp"].(float64))
	ttl := time.Until(time.Unix(expUnix, 0))

	_, err = util.RedisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		if site == "" {
			site = "XXX"
		}

		sessionKey := util.KEY_SESSION + ":" + customerNo + ":" + username

		// Set data
		pipe.HSet(ctx, sessionKey,
			"token:"+strings.ToLower(module), site+":"+tokenString,
			"modules", modulesJson,
		)

		// Set TTL Expire
		pipe.HExpire(ctx, sessionKey, ttl, "token:"+strings.ToLower(module))
		pipe.HExpire(ctx, sessionKey, ttl, "modules")
		pipe.HExpire(ctx, sessionKey, ttl, "company_provider_access")

		// Get IP Address clients
		ipKey := sessionKey + ":ips"
		ipAddress := ctx.ClientIP()
		if ipAddress == "" {
			ipAddress = "unknown"
		}

		// Get existing IP data
		existingData, _ := util.RedisClient.HGet(ctx, ipKey, ipAddress).Result()
		var ipInfo request.IPInfo
		if existingData != "" {
			_ = json.Unmarshal([]byte(existingData), &ipInfo)
		}

		// Update last login
		now := time.Now().Format("2006-01-02 15:04:05")

		// Ensure Modules map is initialized
		if ipInfo.Modules == nil {
			ipInfo.Modules = make(map[string]string)
		}

		// Get User-Agent
		userAgent := ctx.GetHeader("User-Agent")
		if userAgent == "" {
			userAgent = "unknown"
		}

		// Set last login for this module
		ipInfo.Modules[module] = "last_login at " + now + ", user_agent:" + userAgent

		// Save back to Redis
		ipJson, _ := json.Marshal(ipInfo)
		pipe.HSet(ctx, ipKey, ipAddress, ipJson)
		pipe.Expire(ctx, ipKey, ttl)

		return nil
	})
	if err != nil {
		log.Printf("Error storing token in Redis: %v", err)
	}

	// Write response
	util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, resp, "success", 0, 0)
}
