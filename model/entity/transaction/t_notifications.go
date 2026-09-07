package transaction

import (
	"time"
)

type TNotifications struct {
	Id             int64
	CustomerNo     string
	DateLog        time.Time `gorm:"type:date"`
	TimeLog        time.Time
	Username       string
	Name           string
	Reason         string
	Activity       string
	StatusActivity string
	Latitude       float64
	Longitude      float64
	Altitude       float64
	Type           string
	Channel        string
	Category       string
	Title          string
	Content        string
	Event          string
	InsertBy       string
	InsertTime     time.Time
	UpdateBy       string
	UpdateTime     time.Time
}

func (b *TNotifications) TableName() string {
	return "t_notifications"
}
