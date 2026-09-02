package request

import "time"

type TUserGetRequest struct {
	Username   string    `form:"username"`
	Name       string    `form:"name"`
	PhoneNo    string    `form:"phone_no"`
	StatusData string    `form:"status_data"`
	LastTime   time.Time `form:"last_time" time_format:"2006-01-02 15:04:05"`
}

type TUserPostRequest struct {
	Data []TUserPostDetailRequest `json:"data" binding:"required"`
}

type TUserPostDetailRequest struct {
	Type         string `json:"type" binding:"required"`
	Id           int64  `json:"id"`
	IdCustomer   int64  `json:"id_customer" binding:"required"`
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Name         string `json:"name" binding:"required"`
	PhoneNo      string `json:"phone_no" binding:"required"`
	EmailAddress string `json:"email_address" binding:"required"`
	BloodType    string `json:"blood_type" binding:"required"`
	InsertBy     string `json:"insert_by"`
	UpdateBy     string `json:"update_by"`
}
