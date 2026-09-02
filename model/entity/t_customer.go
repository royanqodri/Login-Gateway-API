package entity

import "time"

type TCustomer struct {
	Id          int64
	CustomerNo  string
	Name        string
	Description string
	StatusData  string
	InsertBy    string
	InsertTime  time.Time
	UpdateBy    string
	UpdateTime  time.Time
}

func (b *TCustomer) TableName() string {
	return "t_customer"
}
