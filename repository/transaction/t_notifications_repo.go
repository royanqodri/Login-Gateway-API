package transaction

import (
	"errors"
	"strconv"

	request "github.com/royanqodri/Login-Gateway-API/model/request/transaction"
	"gorm.io/gorm/clause"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	dBConnOperation "github.com/royanqodri/Login-Gateway-API/database"
	entityMaster "github.com/royanqodri/Login-Gateway-API/model/entity"
	entityTransaction "github.com/royanqodri/Login-Gateway-API/model/entity/transaction"
	responseTransaction "github.com/royanqodri/Login-Gateway-API/model/notification"
)

type TNotificationsRepository interface {
	GetByParams(ctx *gin.Context, tx *gorm.DB, dataCustomer entityMaster.TCustomer, request request.TNotificationsGetRequest) (respData []responseTransaction.NotificationMessage, err error)
	Save(ctx *gin.Context, tx *gorm.DB, req []entityTransaction.TNotifications) error
}

type TNotificationsRepositoryImpl struct {
}

func NewTNotificationsRepository() TNotificationsRepository {
	return &TNotificationsRepositoryImpl{}
}

func (repo TNotificationsRepositoryImpl) getTx(tx *gorm.DB) *gorm.DB {
	if tx == nil {
		return dBConnOperation.DBConnOperation
	}
	return tx
}

func (repo TNotificationsRepositoryImpl) GetByParams(ctx *gin.Context, tx *gorm.DB, dataCustomer entityMaster.TCustomer, request request.TNotificationsGetRequest) (respData []responseTransaction.NotificationMessage, err error) {
	result := repo.getTx(tx).
		Select("id, "+strconv.FormatInt(dataCustomer.Id, 10)+" AS id_customer, customer_no, "+
			"to_char(date_log, 'YYYY-MM-DD') AS date_log, shift, shift_sequence, "+
			"equipment_no, equipment_type, equipment_model, "+
			"username, name, fleet, "+
			"CAST(CAST(time_log AS TIMESTAMP(0)) AS VARCHAR) AS time_log, "+
			"state, reason, activity, status_activity, category, "+
			"latitude, longitude, altitude, bearing, type, channel,   "+
			"site, title, content, event, "+
			"insert_by, CAST(CAST(insert_time AS TIMESTAMP(0)) AS VARCHAR) AS insert_time, "+
			"update_by, CAST(CAST(update_time AS TIMESTAMP(0)) AS VARCHAR) AS update_time").
		Table("t_notifications").
		Where("customer_no = ?", dataCustomer.CustomerNo).
		Where("LOWER(type) = ?", "idle_time")

	if request.Site != "" {
		result = result.Where("site = ?", request.Site)
	}

	if request.EquipmentNo != "" {
		result = result.Where("equipment_no = ?", request.EquipmentNo)
	}

	if request.Fleet != "" {
		result = result.Where("fleet = ?", request.Fleet)
	}

	if request.Category != "" {
		result = result.Where("category = ?", request.Category)
	}

	if !request.DateStart.IsZero() && !request.DateEnd.IsZero() {
		result = result.Where("date_log BETWEEN ? AND ?", request.DateStart, request.DateEnd)
	}

	if err = result.Scan(&respData).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []responseTransaction.NotificationMessage{}, nil
		}
		return []responseTransaction.NotificationMessage{}, err
	}

	return respData, nil
}

func (repo *TNotificationsRepositoryImpl) Save(ctx *gin.Context, tx *gorm.DB, req []entityTransaction.TNotifications) error {
	result := repo.getTx(tx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "customer_no"}, {Name: "equipment_no"}, {Name: "time_log"}, {Name: "site"}, {Name: "type"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"update_by", "update_time",
		}),
	}).Create(&req)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
