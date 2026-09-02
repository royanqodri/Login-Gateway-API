package entity

import "time"

type TUser struct {
	Id           int64
	IdCustomer   int64
	Username     string
	Password     string
	Name         string
	PhoneNo      string
	EmailAddress string
	BloodType    string
	StatusData   string
	InsertBy     string
	InsertTime   time.Time
	UpdateBy     string
	UpdateTime   time.Time
}

func (b *TUser) TableName() string {
	return "t_user"
}
