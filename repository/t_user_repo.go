package repository

import (
	"errors"
	"fmt"

	"github.com/royanqodri/Login-Gateway-API/model/entity"
	"github.com/royanqodri/Login-Gateway-API/model/request"
	"github.com/royanqodri/Login-Gateway-API/model/response"
	"github.com/royanqodri/Login-Gateway-API/util"

	"github.com/gin-gonic/gin"
	"github.com/royanqodri/Login-Gateway-API/database"
	"gorm.io/gorm"
)

type TUserRepository interface {
	GetByUsernameOrEmail(ctx *gin.Context, tx *gorm.DB, username string) (response.TUserGetResponse, error)
	GetAll(ctx *gin.Context, tx *gorm.DB, dataCustomer entity.TCustomer, request request.TUserGetRequest, limit, offset int) (respData []response.TUserGetResponse, totalPage int64, totalData int64, err error)
	GetById(ctx *gin.Context, tx *gorm.DB, dataCustomer entity.TCustomer, id int64) (response.TUserGetResponse, error)
	Save(ctx *gin.Context, tx *gorm.DB, req []entity.TUser) error
}

type TUserRepositoryImpl struct {
}

func NewTUserRepository() TUserRepository {
	return &TUserRepositoryImpl{}
}

func (repo TUserRepositoryImpl) getTx(tx *gorm.DB) *gorm.DB {
	if tx == nil {
		return database.DBConnMaster
	}
	return tx
}

func (repo TUserRepositoryImpl) GetByUsernameOrEmail(ctx *gin.Context, tx *gorm.DB, username string) (response.TUserGetResponse, error) {
	data := response.TUserGetResponse{}
	result := repo.getTx(tx).
		Select("t_user.id, " +
			"t_user.id_customer, t_customer.customer_no, " +
			"t_user.username, t_user.password, t_user.name, " +
			"t_user.phone_no, t_user.email_address, t_user.blood_type, " +
			"t_user.status_data, " +
			"t_user.insert_by, CAST(CAST(t_user.insert_time AS TIMESTAMP(0)) AS VARCHAR) AS insert_time, " +
			"t_user.update_by, CAST(CAST(t_user.update_time AS TIMESTAMP(0)) AS VARCHAR) AS update_time").
		Table("t_user").
		Joins("INNER JOIN t_customer ON t_user.id_customer = t_customer.id")

	result = result.Where("(t_user.username = ? OR t_user.email_address = ?)", username, username)
	result.Order("t_user.username ASC")
	result = result.Scan(&data)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return response.TUserGetResponse{}, nil
		}
		return response.TUserGetResponse{}, result.Error
	}

	return data, nil
}

func (repo TUserRepositoryImpl) GetById(ctx *gin.Context, tx *gorm.DB, dataCustomer entity.TCustomer, id int64) (response.TUserGetResponse, error) {
	data := response.TUserGetResponse{}

	result := repo.getTx(tx).
		Select("t_user.id, "+
			"t_user.id_customer, '"+dataCustomer.CustomerNo+"' AS customer_no, "+
			"t_user.username, t_user.password, t_user.name, "+
			"t_user.phone_no, t_user.email_address, t_user.blood_type, t_user.uid_card, "+
			"t_user.status_data, "+
			"t_user.insert_by, CAST(CAST(t_user.insert_time AS TIMESTAMP(0)) AS VARCHAR) AS insert_time, "+
			"t_user.update_by, CAST(CAST(t_user.update_time AS TIMESTAMP(0)) AS VARCHAR) AS update_time").
		Table("t_user").
		Where("t_user.id = ?", id).
		Scan(&data)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return response.TUserGetResponse{}, nil // No data but no error
		}
		return response.TUserGetResponse{}, result.Error // Return error if any
	}

	return data, nil
}

func (repo TUserRepositoryImpl) GetAll(ctx *gin.Context, tx *gorm.DB, dataCustomer entity.TCustomer, request request.TUserGetRequest, limit, offset int) (respData []response.TUserGetResponse, totalPage int64, totalData int64, err error) {
	result := repo.getTx(tx).
		Select("mu.id, "+
			"mu.id_customer, '"+dataCustomer.CustomerNo+"' AS customer_no, "+
			"mu.username, mu.name, mu.phone_no, mu.email_address, mu.blood_type, mu.password, "+
			"mu.status_data AS status_data, mu.uid_card, mu.is_monitor_position, "+
			"mu.insert_by, CAST(CAST(mu.insert_time AS TIMESTAMP(0)) AS VARCHAR) AS insert_time, "+
			"mu.update_by, CAST(CAST(mu.update_time AS TIMESTAMP(0)) AS VARCHAR) AS update_time").
		Table("t_user AS mu").
		Where("mu.id_customer = ?", dataCustomer.Id)

	if request.Username != "" {
		result = result.Where("mu.username = ?", request.Username)
	}

	if request.Name != "" {
		result = result.Where("mu.name LIKE ?", fmt.Sprintf("%%%s%%", request.Name))
	}

	if request.PhoneNo != "" {
		result = result.Where("mu.phone_no = ?", request.PhoneNo)
	}

	statusData := util.GetStatusData(request.StatusData)
	result = result.Where("mu.status_data IN ?", statusData)

	if !request.LastTime.IsZero() {
		result = result.Where("CAST(mu.update_time AS TIMESTAMP(0)) > ?", request.LastTime)
	}

	result = result.Order("mu.username ASC")

	// Count total data
	if err = result.Count(&totalData).Error; err != nil {
		return nil, 0, 0, err
	}

	// Calculate total pages
	if limit > 0 {
		totalPage = (totalData + int64(limit) - 1) / int64(limit)
	} else {
		totalPage = 1
	}

	if limit > 0 {
		result = result.Limit(limit).Offset(offset)
	}

	if err = result.Scan(&respData).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []response.TUserGetResponse{}, totalPage, totalData, nil
		}
		return []response.TUserGetResponse{}, totalPage, totalData, err
	}

	return respData, totalPage, totalData, nil
}

// Save implements TUserRepository.
func (repo TUserRepositoryImpl) Save(ctx *gin.Context, tx *gorm.DB, req []entity.TUser) error {
	for _, user := range req {

		var existing entity.TUser

		query := "id_customer = ? AND (username = ? OR email_address = ? OR phone_no = ?)"
		args := []any{
			user.IdCustomer,
			user.Username,
			user.EmailAddress,
			user.PhoneNo,
		}

		if user.Id > 0 {
			query += " AND id != ?"
			args = append(args, user.Id)
		}

		err := util.CheckDuplicateWithComparator(repo.getTx(tx), &existing, query, args,
			func() []string {
				return util.CompareFields(
					map[string][2]any{
						"id_customer already exists":   {existing.IdCustomer, user.IdCustomer},
						"username already exists":      {existing.Username, user.Username},
						"email_address already exists": {existing.EmailAddress, user.EmailAddress},
						"phone_no already exists":      {existing.PhoneNo, user.PhoneNo},
					},
				)
			},
		)

		if err != nil {
			return err
		}

	}

	result := repo.getTx(tx).Save(&req)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
