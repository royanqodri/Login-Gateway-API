package service

import (
	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
	"github.com/royanqodri/Login-Gateway-API/model/request"
	"github.com/royanqodri/Login-Gateway-API/model/response"
	"github.com/royanqodri/Login-Gateway-API/repository"
)

type SessionService interface {
	GetBySessionData(ctx *gin.Context, sessionDataRequest request.SessionRequest) (respData *response.SessionDataResponse, err error)
	GetUserModules(ctx *gin.Context, request request.RedisRequest) ([]response.TUserModuleGetResponse, error)
}

type SessionServiceImpl struct {
	mstUserModuleRepo repository.TUserModuleRepository
	mstUserRepo       repository.TUserRepository
	redisClient       *redis.Client
}

func NewSessionService(
	mstUserModuleRepo repository.TUserModuleRepository,
	mstUserRepo repository.TUserRepository,
	redisClient *redis.Client) SessionService {
	return &SessionServiceImpl{
		mstUserModuleRepo: mstUserModuleRepo,
		mstUserRepo:       mstUserRepo,
		redisClient:       redisClient,
	}
}

func (service SessionServiceImpl) GetBySessionData(ctx *gin.Context, sessionDataRequest request.SessionRequest) (respData *response.SessionDataResponse, err error) {
	// Get detail data user by customer_no and username
	respDataUser, err := service.mstUserRepo.GetByUsernameOrEmail(ctx, nil, sessionDataRequest.Username)
	if err != nil {

		return nil, err
	}

	// Get Data Permissions Module by customer_no, username, site and system code (menu)
	rawData, err := service.mstUserModuleRepo.GetByUsernameAndSiteAndMenu(ctx, nil, sessionDataRequest.CustomerNo, sessionDataRequest.Username, sessionDataRequest.Site, sessionDataRequest.Menu)
	if err != nil {
		return nil, err
	}

	userData := &response.SessionDataResponse{
		IdCustomer:   respDataUser.IdCustomer,
		CustomerNo:   respDataUser.CustomerNo,
		IdUser:       respDataUser.Id,
		Username:     respDataUser.Username,
		Name:         respDataUser.Name,
		PhoneNo:      respDataUser.PhoneNo,
		EmailAddress: respDataUser.EmailAddress,
		BloodType:    respDataUser.BloodType,
		DataModules:  []response.TUserModuleResponse{},
	}

	// Populate modules for the user
	for _, item := range rawData {
		userData.DataModules = append(userData.DataModules, response.TUserModuleResponse{
			Module:       item.Module,
			CreateAccess: item.CreateAccess,
			ReadAccess:   item.ReadAccess,
			UpdateAccess: item.UpdateAccess,
			DeleteAccess: item.DeleteAccess,
			ExportAccess: item.ExportAccess,
		})
	}

	return userData, nil
}

func (service SessionServiceImpl) GetUserModules(ctx *gin.Context, request request.RedisRequest) ([]response.TUserModuleGetResponse, error) {
	// Fetch user permissions from the database based on the provided parameters
	userModules, err := service.mstUserModuleRepo.GetByUsernameAndSiteAndMenu(ctx, nil, request.CustomerNo, request.Username, request.Site, request.Menu)
	if err != nil {
		return nil, err
	}

	return userModules, nil
}
