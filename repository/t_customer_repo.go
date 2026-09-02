package repository

import (
	"errors"

	"github.com/royanqodri/Login-Gateway-API/database"
	entityMaster "github.com/royanqodri/Login-Gateway-API/model/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TCustomerRepository interface {
	GetById(ctx *gin.Context, tx *gorm.DB, id int64) (entityMaster.TCustomer, error)
	GetByCustomerNo(ctx *gin.Context, tx *gorm.DB, customerNo string) (entityMaster.TCustomer, error)
	GetAll(ctx *gin.Context, tx *gorm.DB, customerNo string, req entityMaster.TCustomer) ([]entityMaster.TCustomer, error)
}

type TCustomerRepositoryImpl struct {
}

func NewTCustomerRepository() TCustomerRepository {
	return &TCustomerRepositoryImpl{}
}

func (repo TCustomerRepositoryImpl) getTx(tx *gorm.DB) *gorm.DB {
	if tx == nil {
		return database.DBConnMaster
	}
	return tx
}

func (repo TCustomerRepositoryImpl) GetById(ctx *gin.Context, tx *gorm.DB, id int64) (entityMaster.TCustomer, error) {
	data := entityMaster.TCustomer{}
	result := repo.getTx(tx).First(&data, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {

			return entityMaster.TCustomer{}, nil
		}
		return entityMaster.TCustomer{}, result.Error
	}

	return data, nil
}

func (repo TCustomerRepositoryImpl) GetByCustomerNo(ctx *gin.Context, tx *gorm.DB, customerNo string) (entityMaster.TCustomer, error) {
	data := entityMaster.TCustomer{}
	result := repo.getTx(tx).Where(entityMaster.TCustomer{CustomerNo: customerNo}).First(&data)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entityMaster.TCustomer{}, nil
		}
		return entityMaster.TCustomer{}, result.Error

	}

	return data, nil
}

func (repo TCustomerRepositoryImpl) GetAll(ctx *gin.Context, tx *gorm.DB, customerNo string, req entityMaster.TCustomer) ([]entityMaster.TCustomer, error) {
	data := []entityMaster.TCustomer{}
	result := repo.getTx(tx)

	if customerNo != "" {
		result = result.Where(entityMaster.TCustomer{CustomerNo: req.CustomerNo})
	}

	if req.Name != "" {
		result = result.Where(entityMaster.TCustomer{Name: req.Name})

	}

	result = result.Find(&data)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {

			return []entityMaster.TCustomer{}, nil
		}
		return []entityMaster.TCustomer{}, result.Error

	}

	return data, nil
}
