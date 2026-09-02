package service

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/royanqodri/Login-Gateway-API/model/request"
	"github.com/royanqodri/Login-Gateway-API/model/response"
	"github.com/royanqodri/Login-Gateway-API/repository"
	"github.com/royanqodri/Login-Gateway-API/util"
	"github.com/royanqodri/Login-Gateway-API/util/constants"

	"github.com/gin-gonic/gin"
)

type LoginService interface {
	Login(ctx *gin.Context, request request.LoginRequest) (response.LoginResponse, error)
	// OnboardGenerateToken(ctx *gin.Context, request request.OnboardRequest) (response.OnboardResponse, error)
	LoginWithGoogle(ctx *gin.Context, idToken string) (response.LoginResponse, error)
	LoginWithFacebook(ctx *gin.Context, accessToken string) (response.LoginResponse, error)
	StoreSessionInRedis(ctx *gin.Context, customerNo string, username string, site string, module string, token string, userModules []response.TUserModuleGetResponse, ttl time.Duration) error
	GetUserModules(ctx *gin.Context, request request.RedisRequest) ([]response.TUserModuleGetResponse, error)
}

type LoginServiceImpl struct {
	mstUserRepo       repository.TUserRepository
	mstUserModuleRepo repository.TUserModuleRepository
	mstCustomerRepo   repository.TCustomerRepository
}

func NewLoginService(
	mstUserRepo repository.TUserRepository,
	mstUserModuleRepo repository.TUserModuleRepository,

	mstCustomerRepo repository.TCustomerRepository,
) LoginService {
	return &LoginServiceImpl{
		mstUserRepo:       mstUserRepo,
		mstUserModuleRepo: mstUserModuleRepo,
		mstCustomerRepo:   mstCustomerRepo,
	}
}

func (service LoginServiceImpl) Login(ctx *gin.Context, request request.LoginRequest) (response.LoginResponse, error) {
	// get data user
	user, err := service.mstUserRepo.GetByUsernameOrEmail(ctx, nil, request.Username)
	if err != nil {
		return response.LoginResponse{}, err
	}

	if user.Id == 0 {
		return response.LoginResponse{}, util.NewErrorMessage(400, constants.USER_NOT_FOUND, nil)
	}

	// compare password
	isPasswordMatched := util.ComparePasswordWithMD5(request.Password, user.Password)
	if !isPasswordMatched {
		return response.LoginResponse{}, util.NewErrorMessage(400, constants.PASSWORD_IS_WRONG, nil)
	}

	// generate jwt - authorization token
	token, err := util.GenerateToken(user.Id, user.Username, request.Site, request.Menu, user.CustomerNo, user.IdCustomer)
	if err != nil {
		return response.LoginResponse{}, err
	}

	loginResponse := response.LoginResponse{
		IdCustomer: user.IdCustomer,
		CustomerNo: user.CustomerNo,
		Username:   user.Username,
		Name:       user.Name,
		Token:      token,
	}

	return loginResponse, nil
}

// func (service *LoginServiceImpl) OnboardGenerateToken(ctx *gin.Context, request request.OnboardRequest) (response.OnboardResponse, error) {
// 	// Get Data Device by Serial and License
// 	device, err := service.mstDeviceRepo.GetBySerial(ctx, nil, request)
// 	if err != nil {
// 		return response.OnboardResponse{}, err
// 	}

// 	if device.Id == 0 {
// 		return response.OnboardResponse{}, util.NewErrorMessage(400, constants.DEVICE_NOT_FOUND, nil)
// 	}

// 	customerNo := device.CustomerNo
// 	customID := device.IdCustomer

// 	// generate jwt - authorization token
// 	token, err := util.GenerateTokenOnboard(request.Menu, customerNo, customID)
// 	if err != nil {
// 		return response.OnboardResponse{}, err
// 	}

// 	onboardResponse := response.OnboardResponse{
// 		IdCustomer: customID,
// 		CustomerNo: customerNo,
// 		Token:      token,
// 	}

// 	return onboardResponse, nil
// }

func (service LoginServiceImpl) StoreSessionInRedis(ctx *gin.Context, customerNo string, username string, site string, module string, token string, userModules []response.TUserModuleGetResponse, ttl time.Duration) error {
	// Serialize the user modules (permissions) into JSON

	modulesJson, err := json.Marshal(userModules)
	if err != nil {
		log.Printf("Error marshaling user modules: %v", err)
		return err
	}

	_, err = util.RedisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		if site == "" {
			site = "XXX"
		}

		sessionKey := util.KEY_SESSION + ":" + customerNo + ":" + username

		// Set data
		pipe.HSet(ctx, sessionKey,
			"token:"+strings.ToLower(module), site+":"+token,
			"modules", modulesJson,
		)

		// Set TTL Expire
		pipe.HExpire(ctx, sessionKey, ttl, "token:"+strings.ToLower(module))
		pipe.HExpire(ctx, sessionKey, ttl, "modules")
		pipe.HExpire(ctx, sessionKey, ttl, "company_provider_access")

		return nil
	})
	if err != nil {
		log.Printf("Error storing token in Redis: %v", err)
		return err
	}

	return nil
}

func (service LoginServiceImpl) GetUserModules(ctx *gin.Context, request request.RedisRequest) ([]response.TUserModuleGetResponse, error) {
	// Fetch user permissions from the database based on the provided parameters
	userModules, err := service.mstUserModuleRepo.GetByUsernameAndSiteAndMenu(ctx, nil, request.CustomerNo, request.Username, request.Site, request.Menu)
	if err != nil {
		return nil, err
	}

	return userModules, nil
}

func (service LoginServiceImpl) LoginWithGoogle(ctx *gin.Context, idToken string) (response.LoginResponse, error) {
	// 1. Verifikasi ID Token dari Google
	googleUser, err := util.VerifyGoogleIDToken(ctx.Request.Context(), idToken)
	if err != nil {
		return response.LoginResponse{}, util.NewErrorMessage(401, "INVALID_GOOGLE_TOKEN", err)
	}

	// GET by email
	user, err := service.mstUserRepo.GetByUsernameOrEmail(ctx, nil, googleUser.Email)
	if err != nil {
		return response.LoginResponse{}, err
	}

	if user.Id == 0 {
		return response.LoginResponse{}, util.NewErrorMessage(400, constants.USER_NOT_FOUND, nil)

		// To do register google
		/*
			user, err = service.mstUserRepo.CreateFromGoogle(ctx, googleUser)
			if err != nil {
				return response.LoginResponse{}, err
			}
		*/
	}

	token, err := util.GenerateToken(user.Id, user.Username, "", "", user.CustomerNo, user.IdCustomer)
	if err != nil {
		return response.LoginResponse{}, err
	}

	loginResponse := response.LoginResponse{
		IdCustomer: user.IdCustomer,
		CustomerNo: user.CustomerNo,
		Username:   user.Username,
		Name:       user.Name,
		Token:      token,
	}

	return loginResponse, nil
}

func (service LoginServiceImpl) LoginWithFacebook(ctx *gin.Context, accessToken string) (response.LoginResponse, error) {

	facebookUser, err := util.VerifyFacebookAccessToken(ctx.Request.Context(), accessToken)
	if err != nil {
		return response.LoginResponse{}, util.NewErrorMessage(401, "INVALID_FACEBOOK_TOKEN", err)
	}

	// GET by email
	user, err := service.mstUserRepo.GetByUsernameOrEmail(ctx, nil, facebookUser.Email)
	if err != nil {
		return response.LoginResponse{}, err
	}

	if user.Id == 0 {
		return response.LoginResponse{}, util.NewErrorMessage(400, constants.USER_NOT_FOUND, nil)

		// To do register facebook
		/*
			user, err = service.mstUserRepo.CreateFromFacebook(ctx, facebookUser)
			if err != nil {
				return response.LoginResponse{}, err
			}
		*/
	}

	token, err := util.GenerateToken(user.Id, user.Username, "", "", user.CustomerNo, user.IdCustomer)
	if err != nil {
		return response.LoginResponse{}, err
	}

	loginResponse := response.LoginResponse{
		IdCustomer: user.IdCustomer,
		CustomerNo: user.CustomerNo,
		Username:   user.Username,
		Name:       user.Name,
		Token:      token,
	}

	return loginResponse, nil
}
