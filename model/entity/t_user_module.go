package entity

import "time"

type TUserModule struct {
	Id           int64
	IdCustomer   int64
	IdUser       int64
	IdModule     int64
	CreateAccess bool
	ReadAccess   bool
	UpdateAccess bool
	DeleteAccess bool
	ExportAccess bool
	StatusData   string
	InsertBy     string
	InsertTime   time.Time
	UpdateBy     string
	UpdateTime   time.Time
}

func (b *TUserModule) TableName() string {
	return "t_user_module"
}
