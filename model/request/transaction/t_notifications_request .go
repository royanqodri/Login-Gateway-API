package transaction

import (
	"time"
)

type TNotificationsGetRequest struct {
	Site        string    `form:"site"`
	Category    string    `form:"category"`
	EquipmentNo string    `form:"equipment_no"`
	Fleet       string    `form:"fleet"`
	Type        string    `form:"type"`
	DateStart   time.Time `form:"date_start" time_format:"2006-01-02" binding:"required"`
	DateEnd     time.Time `form:"date_end" time_format:"2006-01-02" binding:"required"`
}

type TNotificationsPostRequest struct {
	Data []TNotificationsDetailRequest `json:"data" binding:"required"`
}

type TNotificationsDetailRequest struct {
	Id             int64   `json:"id"`
	CustomerNo     string  `json:"customer_no"`
	DateLog        string  `json:"date_log" time_format:"2006-01-02"`
	TimeLog        string  `json:"time_log" time_format:"2006-01-02 15:04:05"`
	Shift          string  `json:"shift"`
	ShiftSequence  int64   `json:"shift_sequence"`
	EquipmentNo    string  `json:"equipment_no"`
	EquipmentType  string  `json:"equipment_type"`
	EquipmentModel string  `json:"equipment_model"`
	Username       string  `json:"username"`
	Name           string  `json:"name"`
	Fleet          string  `json:"fleet"`
	State          string  `json:"state"`
	Reason         string  `json:"reason"`
	Activity       string  `json:"activity"`
	StatusActivity string  `json:"status_activity"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Altitude       float64 `json:"altitude"`
	Bearing        float64 `json:"bearing"`
	Channel        string  `json:"channel"`
	Category       string  `json:"category"`
	Title          string  `json:"title"`
	Content        string  `json:"content"`
	Type           string  `json:"type"`
	Event          string  `json:"event"`
	Site           string  `json:"site"`
	InsertBy       string  `json:"insert_by"`
	UpdateBy       string  `json:"update_by"`
}
