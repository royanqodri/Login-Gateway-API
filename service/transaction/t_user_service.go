package transaction

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/royanqodri/Login-Gateway-API/database"
	"github.com/royanqodri/Login-Gateway-API/model/entity"
	"github.com/royanqodri/Login-Gateway-API/model/request"
	"github.com/royanqodri/Login-Gateway-API/model/response"
	"github.com/royanqodri/Login-Gateway-API/repository"
	"github.com/royanqodri/Login-Gateway-API/util"
	"github.com/royanqodri/Login-Gateway-API/util/constants"
)

type TUserService interface {
	GetAll(ctx *gin.Context, dataCustomer entity.TCustomer, request request.TUserGetRequest) (respData []response.TUserGetResponse, totalPage int64, totalData int64, err error)
	Post(ctx *gin.Context, dataCustomer entity.TCustomer, request request.TUserPostRequest) (err error)
}

type TUserServiceImpl struct {
	tUserRepo     repository.TUserRepository
	tCustomerRepo repository.TCustomerRepository
}

func NewTUserService(tUserRepo repository.TUserRepository, tCustomerRepo repository.TCustomerRepository) TUserService {
	return &TUserServiceImpl{
		tUserRepo:     tUserRepo,
		tCustomerRepo: tCustomerRepo,
	}
}

func (service TUserServiceImpl) GetAll(ctx *gin.Context, dataCustomer entity.TCustomer, request request.TUserGetRequest) (respData []response.TUserGetResponse, totalPage int64, totalData int64, err error) {

	// Parse pagination parameters
	page, pageProvided := ctx.GetQuery("page_now")
	limit, limitProvided := ctx.GetQuery("limit")

	// Default values
	var pageNum int64 = 1
	var limitNum int64 = 0
	var offset int64 = 0

	// Determine pagination logic
	if pageProvided || limitProvided {
		if pageProvided {
			pageNum, err = strconv.ParseInt(page, 10, 64)
			if err != nil || pageNum < 1 {
				pageNum = 1
			}
		}

		if limitProvided {
			limitNum, err = strconv.ParseInt(limit, 10, 64)
			if err != nil || limitNum < 1 {
				limitNum = 10
			}
		} else if pageProvided {
			limitNum = 10 // Default 10 BE  when only `page_now` is provided
		}

		offset = (pageNum - 1) * limitNum
	} else {

		limitNum = 0
	}

	respData, _, totalData, err = service.tUserRepo.GetAll(ctx, nil, dataCustomer, request, int(limitNum), int(offset))
	if err != nil {

		return []response.TUserGetResponse{}, 0, 0, fmt.Errorf("failed to get paginated data: %w", err)
	}
	// Calculate total pages
	if limitNum > 0 && totalData > 0 {
		totalPage = (totalData + limitNum - 1) / limitNum
	} else {
		totalPage = 1 // Default to 1 page if no pagination
	}

	// Return early if no data
	if totalData == 0 {
		return []response.TUserGetResponse{}, totalPage, totalData, nil
	}
	return respData, totalPage, totalData, nil
}

func (service TUserServiceImpl) Post(ctx *gin.Context, dataCustomer entity.TCustomer, request request.TUserPostRequest) (err error) {
	tx := database.DBConnMaster.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Convert post request data to entities
	dataReq := make([]entity.TUser, len(request.Data))
	for i, data := range request.Data {
		if data.Type == constants.STATUS_DATA_INSERT {
			// INSERT Operation
			passwordHash, err := util.HashPasswordWithMD5(data.Password)
			if err != nil {
				return err
			}

			dataReq[i] = entity.TUser{
				IdCustomer:   data.IdCustomer,
				Username:     strings.TrimSpace(data.Username),
				Password:     passwordHash,
				Name:         strings.TrimSpace(data.Name),
				PhoneNo:      strings.TrimSpace(data.PhoneNo),
				EmailAddress: strings.TrimSpace(data.EmailAddress),
				BloodType:    strings.TrimSpace(data.BloodType),
				StatusData:   strings.TrimSpace(data.Type),
				InsertBy:     strings.TrimSpace(data.InsertBy),
				InsertTime:   util.GetTimeNowByLoc(),
				UpdateBy:     strings.TrimSpace(data.UpdateBy),
				UpdateTime:   util.GetTimeNowByLoc(),
			}
		} else if data.Type == constants.STATUS_DATA_UPDATE || data.Type == constants.STATUS_DATA_DELETE {
			// For update or delete, retrieve the existing record to preserve InsertTime

			existingData, err := service.tUserRepo.GetById(ctx, tx, dataCustomer, data.Id)
			if err != nil {
				return err
			}

			parsedInsertTime, err := time.Parse("2006-01-02 15:04:05", existingData.InsertTime)
			if err != nil {
				return err
			}

			passwordHash := existingData.Password
			if data.Password != "" && strings.TrimSpace(data.Password) != "" {
				isPasswordMatched := util.ComparePasswordWithMD5(data.Password, existingData.Password)
				if !isPasswordMatched {
					passwordUpdate, err := util.HashPasswordWithMD5(data.Password)
					if err != nil {
						return err
					}
					passwordHash = passwordUpdate
				}
			}

			dataReq[i] = entity.TUser{
				Id:           data.Id,
				IdCustomer:   data.IdCustomer,
				Username:     strings.TrimSpace(data.Username),
				Password:     passwordHash,
				Name:         strings.TrimSpace(data.Name),
				PhoneNo:      strings.TrimSpace(data.PhoneNo),
				EmailAddress: strings.TrimSpace(data.EmailAddress),
				BloodType:    strings.TrimSpace(data.BloodType),
				StatusData:   strings.TrimSpace(data.Type),
				InsertBy:     existingData.InsertBy,
				InsertTime:   parsedInsertTime,
				UpdateBy:     strings.TrimSpace(data.UpdateBy),
				UpdateTime:   util.GetTimeNowByLoc(),
			}
		}
	}

	// Save the data to the database
	if len(dataReq) > 0 {
		if err = service.tUserRepo.Save(ctx, tx, dataReq); err != nil {
			return err
		}
	}

	// Commit the transaction
	if err = tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
