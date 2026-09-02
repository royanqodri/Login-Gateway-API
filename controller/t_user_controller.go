package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/royanqodri/Login-Gateway-API/model/entity"
	"github.com/royanqodri/Login-Gateway-API/model/request"
	_ "github.com/royanqodri/Login-Gateway-API/model/response"
	"github.com/royanqodri/Login-Gateway-API/service/transaction"
	"github.com/royanqodri/Login-Gateway-API/util"
	"github.com/royanqodri/Login-Gateway-API/util/logging"
	"github.com/sirupsen/logrus"
)

type TUserController interface {
	GetAll(ctx *gin.Context)
	Post(ctx *gin.Context)
}

type TUserControllerImpl struct {
	tUserService transaction.TUserService
}

func NewTUserController(tUserService transaction.TUserService) TUserController {
	return &TUserControllerImpl{tUserService: tUserService}
}

// GetAll godoc
// @Summary Get Master User
// @Description Get list of User by filter
// @Tags Master User
// @Param customer_no path string true "Customer No"
// @Param site query string false "Site"
// @Param department query string false "Department"
// @Param section query string false "Section"
// @Param position query string false "Position"
// @Param username query string false "Username"
// @Param name query string false "Name"
// @Param crew query string false "Crew"
// @Param phone_no query string false "Phone Number"
// @Param status_data query string false "Status Data"
// @Param source query string true "Source"
// @Param last_time query string false "Last Update Time" Format("2006-01-02 15:04:05")
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} response.TUserMainResponse
// @Failure 400 {object} response.MainResponse
// @Failure 500 {object} response.MainResponse
// @Router /user [get]
func (controller *TUserControllerImpl) GetAll(ctx *gin.Context) {
	dataCustomer := entity.TCustomer{
		Id:         ctx.GetInt64("customer_id"),
		CustomerNo: ctx.Param("customer_no"),
	}

	req := request.TUserGetRequest{}

	// Validate request payload
	if err := ctx.ShouldBind(&req); err != nil {
		logging.LogWithFields(logging.ERROR, logging.ERROR, logrus.Fields{
			"endpoint": ctx.Request.URL,
			"method":   ctx.Request.Method,
			"message":  err,
		})
		util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	// Get data from transaction
	resp, totalPage, totalData, err := controller.tUserService.GetAll(ctx, dataCustomer, req)
	if err != nil {
		logging.LogWithFields(logging.ERROR, logging.ERROR, logrus.Fields{
			"endpoint": ctx.Request.URL,
			"method":   ctx.Request.Method,
			"message":  err,
		})
		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	// Write response
	util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, resp, "success", totalPage, totalData)
}

// Post godoc
// @Summary Create Master User
// @Description Create one or more master User entries
// @Tags Master User
// @Param customer_no path string true "Customer No"
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param request body request.TUserPostRequest true "User Post Payload"
// @Success 200 {object} response.MainResponse
// @Failure 400 {object} response.MainResponse
// @Failure 409 {object} response.MainResponse
// @Failure 500 {object} response.MainResponse
// @Router /user [post]
func (controller *TUserControllerImpl) Post(ctx *gin.Context) {
	dataCustomer := entity.TCustomer{
		Id:         ctx.GetInt64("customer_id"),
		CustomerNo: ctx.Param("customer_no"),
	}

	// Validate request payload
	req := request.TUserPostRequest{}

	err := ctx.ShouldBind(&req)
	if err != nil {
		logging.LogWithFields(logging.ERROR, logging.ERROR, logrus.Fields{

			"endpoint": ctx.Request.URL,
			"method":   ctx.Request.Method,
			"message":  err,
		})
		util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	err = controller.tUserService.Post(ctx, dataCustomer, req)
	if err != nil {
		logging.LogWithFields(logging.ERROR, logging.ERROR, logrus.Fields{

			"endpoint": ctx.Request.URL,
			"method":   ctx.Request.Method,
			"message":  err,
		})
		var dupErr *util.DuplicateError
		if errors.As(err, &dupErr) {
			util.WriteDuplicateResponse(ctx, dupErr.Fields)
			return
		}

		util.WriteResponse(ctx, util.JSON, http.StatusInternalServerError, []any{}, []any{}, fmt.Sprintf("%v", err), 0, 0)
		return
	}

	// Write response
	util.WriteResponse(ctx, util.JSON, http.StatusOK, []any{}, []any{}, "success", 0, 0)
}
