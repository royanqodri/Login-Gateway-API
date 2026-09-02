package repository

import (
	"errors"

	"github.com/royanqodri/Login-Gateway-API/model/response"
	"github.com/royanqodri/Login-Gateway-API/util"

	"github.com/gin-gonic/gin"
	"github.com/royanqodri/Login-Gateway-API/database"
	"gorm.io/gorm"
)

type TUserModuleRepository interface {
	GetByUsernameAndSiteAndMenu(ctx *gin.Context, tx *gorm.DB, customerNo string, username string, site string, menu string) ([]response.TUserModuleGetResponse, error)
}

type TUserModuleRepositoryImpl struct {
}

func NewTUserModuleRepository() TUserModuleRepository {
	return &TUserModuleRepositoryImpl{}
}

func (repo TUserModuleRepositoryImpl) getTx(tx *gorm.DB) *gorm.DB {
	if tx == nil {
		return database.DBConnMaster
	}
	return tx
}

func (repo TUserModuleRepositoryImpl) GetByUsernameAndSiteAndMenu(ctx *gin.Context, tx *gorm.DB, customerNo string, username string, site string, menu string) ([]response.TUserModuleGetResponse, error) {
	data := []response.TUserModuleGetResponse{}
	result := repo.getTx(tx).
		Select("t_user_module.id, " +
			"t_user_module.id_customer, " +
			"t_user_module.id_user, " +
			"t_user_module.id_module, " +
			"t_module.code as module, " +
			"t_user_module.create_access, " +
			"t_user_module.read_access, " +
			"t_user_module.update_access, " +
			"t_user_module.delete_access, " +
			"t_user_module.export_access, " +
			"t_user_module.status_data, " +
			"t_user_module.insert_by, " +
			"t_user_module.insert_time, " +
			"t_user_module.update_by, " +
			"t_user_module.update_time").
		Table("t_user_module")

	if customerNo != "" {
		result = result.Joins("INNER JOIN t_customer ON t_user_module.id_customer = t_customer.id").
			Where("t_customer.customer_no = ?", customerNo)
	}

	result = result.
		Joins("INNER JOIN t_user ON t_user_module.id_user = t_user.id").
		Joins("INNER JOIN t_module ON t_user_module.id_module = t_module.id").
		Joins("INNER JOIN t_system ON t_module.id_system = t_system.id").
		Where("t_user_module.status_data != ?", "D").
		Order("t_user.username ASC, t_module.code ASC")

	if username != "" {
		result = result.Where("t_user.username = ?", username)
	}

	statusData := util.GetStatusData("A")
	result = result.Where("t_module.status_data IN ?", statusData)

	if menu != "" {
		result = result.Where("t_system.code = ?", menu)
	}

	result.Scan(&data)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []response.TUserModuleGetResponse{}, nil
		}
		return []response.TUserModuleGetResponse{}, result.Error
	}

	return data, nil
}
