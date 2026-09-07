package transaction

import (
	"github.com/royanqodri/Login-Gateway-API/database"
	entityMaster "github.com/royanqodri/Login-Gateway-API/model/entity"
	entityTransaction "github.com/royanqodri/Login-Gateway-API/model/entity/transaction"
	responseTransaction "github.com/royanqodri/Login-Gateway-API/model/notification"
	requestTransaction "github.com/royanqodri/Login-Gateway-API/model/request/transaction"
	repositoryMaster "github.com/royanqodri/Login-Gateway-API/repository"
	repositoryTransaction "github.com/royanqodri/Login-Gateway-API/repository/transaction"
	"github.com/royanqodri/Login-Gateway-API/util"

	"github.com/gin-gonic/gin"
)

type TNotificationsService interface {
	GetByParams(ctx *gin.Context, dataCustomer entityMaster.TCustomer, request requestTransaction.TNotificationsGetRequest) (respData []responseTransaction.NotificationMessage, err error)
	Post(ctx *gin.Context, dataCustomer entityMaster.TCustomer, request requestTransaction.TNotificationsPostRequest) (err error)
}

type TNotificationsServiceImpl struct {
	tNotificationsRepo repositoryTransaction.TNotificationsRepository
	TCustomerRepo      repositoryMaster.TCustomerRepository
}

func NewTNotificationsService(tNotificationsRepo repositoryTransaction.TNotificationsRepository, TCustomerRepo repositoryMaster.TCustomerRepository) TNotificationsService {
	return &TNotificationsServiceImpl{
		tNotificationsRepo: tNotificationsRepo,
		TCustomerRepo:      TCustomerRepo,
	}
}

func (service TNotificationsServiceImpl) GetByParams(ctx *gin.Context, dataCustomer entityMaster.TCustomer, request requestTransaction.TNotificationsGetRequest) (respData []responseTransaction.NotificationMessage, err error) {
	// Get Data from DB
	respData, err = service.tNotificationsRepo.GetByParams(ctx, nil, dataCustomer, request)
	if err != nil {
		return []responseTransaction.NotificationMessage{}, err
	}

	return respData, nil
}

func (service TNotificationsServiceImpl) Post(ctx *gin.Context, dataCustomer entityMaster.TCustomer, request requestTransaction.TNotificationsPostRequest) (err error) {
	tx := database.DBConnOperation.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	dataReq := make([]entityTransaction.TNotifications, len(request.Data))
	for i, data := range request.Data {
		dataReq[i] = entityTransaction.TNotifications{
			CustomerNo:     dataCustomer.CustomerNo,
			DateLog:        util.GetFormattedDate(data.DateLog),
			TimeLog:        util.GetFormattedDateTime(data.TimeLog),
			Username:       data.Username,
			Name:           data.Name,
			Reason:         data.Reason,
			Activity:       data.Activity,
			StatusActivity: data.StatusActivity,
			Latitude:       util.RoundToPrecision(data.Latitude, 13),
			Longitude:      util.RoundToPrecision(data.Longitude, 13),
			Altitude:       util.RoundToPrecision(data.Altitude, 13),
			Type:           data.Type,
			Channel:        data.Channel,
			Category:       data.Category,
			Title:          data.Title,
			Content:        data.Content,
			Event:          data.Event,
			InsertBy:       ctx.GetString("user"),
			InsertTime:     util.GetTimeNowByLoc(),
			UpdateBy:       ctx.GetString("user"),
			UpdateTime:     util.GetTimeNowByLoc(),
		}

		if ctx.GetString("user") == "" {
			dataReq[i].InsertBy = data.InsertBy
			dataReq[i].UpdateBy = data.UpdateBy
		}
	}

	if len(dataReq) > 0 {
		if err = service.tNotificationsRepo.Save(ctx, tx, dataReq); err != nil {
			return err
		}
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
