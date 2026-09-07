package controller

import (
	"fmt"
	"time"

	"github.com/royanqodri/Login-Gateway-API/config"
	"github.com/royanqodri/Login-Gateway-API/model/request"
	_ "github.com/royanqodri/Login-Gateway-API/model/response"
	service "github.com/royanqodri/Login-Gateway-API/service"
	"github.com/royanqodri/Login-Gateway-API/util"

	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginController interface {
	Login(ctx *gin.Context)
	LoginWithRedis(ctx *gin.Context)
	LoginWithGoogle(ctx *gin.Context)
	LoginWithFacebook(ctx *gin.Context)
	LoginWithApple(ctx *gin.Context)
}

type LoginControllerImpl struct {
	loginService service.LoginService
}

func NewLoginController(loginService service.LoginService) LoginController {
	return &LoginControllerImpl{loginService: loginService}
}

func (controller *LoginControllerImpl) Login(ctx *gin.Context) {
	// validate request payload
	req := request.LoginRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	resp, err := controller.loginService.Login(ctx, req)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	// write response
	util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, resp, "success", 0, 0)
}

// LoginWithRedis godoc
// @Summary Login With Redis
// @Description Create Authenticate user, store in redis, and return a JWT token
// @Tags Login
// @Accept  json
// @Produce  json
// @Param request body request.LoginRequest true "Login Post Payload"
// @Success 200 {object} response.LoginMainResponse
// @Failure 400 {object} response.MainResponse
// @Failure 409 {object} response.MainResponse
// @Failure 500 {object} response.MainResponse
// @Router /mintegra/login [post]
func (controller LoginControllerImpl) LoginWithRedis(ctx *gin.Context) {
	// Bind login request
	var req request.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, fmt.Sprintf("Invalid request payload: %v", err), 0, 0)
		return
	}

	// Authenticate and generate token
	resp, err := controller.loginService.Login(ctx, req)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("Login failed: %v", err), 0, 0)
		return
	}

	// Build Redis request payload
	req2 := request.RedisRequest{
		Username:   resp.Username,
		CustomerID: resp.IdCustomer,
		CustomerNo: resp.CustomerNo,
		Menu:       "",
	}

	// Fetch user modules
	userModules, err := controller.loginService.GetUserModules(ctx, req2)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("Error fetching user permissions: %v", err), 0, 0)
		return
	}

	// Store token in Redis
	expUnix := util.GetTimeNowByLoc().Add(time.Duration(config.Get().JWT.Expire) * time.Second).Unix()
	ttl := time.Until(time.Unix(expUnix, 0))
	err = controller.loginService.StoreSessionInRedis(ctx, resp.CustomerNo, resp.Username, req.Site, req.Menu, resp.Token, userModules, ttl)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("Error storing token in Redis: %v", err), 0, 0)
		return
	}

	// Return token in response
	util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, resp, "success", 0, 0)
}

func (controller *LoginControllerImpl) LoginWithGoogle(ctx *gin.Context) {
	// validate request payload
	req := request.GoogleLoginRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	resp, err := controller.loginService.LoginWithGoogle(ctx, req.IdToken)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	// write response
	util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, resp, "success", 0, 0)
}

func (controller *LoginControllerImpl) LoginWithFacebook(ctx *gin.Context) {
	// validate request payload
	req := request.FacebookLoginRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	resp, err := controller.loginService.LoginWithFacebook(ctx, req.AccessToken)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	// write response
	util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, resp, "success", 0, 0)
}

func (controller *LoginControllerImpl) LoginWithApple(ctx *gin.Context) {
	// validate request payload
	req := request.AppleLoginRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	resp, err := controller.loginService.LoginWithApple(ctx, req.IdToken)
	if err != nil {
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	// write response
	util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, resp, "success", 0, 0)
}
