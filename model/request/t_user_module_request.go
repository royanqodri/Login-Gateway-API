package request

import "time"

type TUserModuleGetRequest struct {
	Username string    `form:"username"`
	Menu     string    `form:"menu"`
	LastTime time.Time `form:"last_time" time_format:"2006-01-02 15:04:05"`
}

type TUserModulePostRequest struct {
	Data []TUserModulePostDetailRequest `json:"data" binding:"required"`
}

type TUserModulePostDetailRequest struct {
	Type         string `json:"type" binding:"required"`
	Id           int64  `json:"id"`
	IdCustomer   int64  `json:"id_customer" binding:"required"`
	IdUser       int64  `json:"id_user" binding:"required"`
	IdModule     int64  `json:"id_module" binding:"required"`
	CreateAccess bool   `json:"create_access" binding:"required"`
	ReadAccess   bool   `json:"read_access" binding:"required"`
	UpdateAccess bool   `json:"update_access" binding:"required"`
	DeleteAccess bool   `json:"delete_access" binding:"required"`
	ExportAccess bool   `json:"export_access" binding:"required"`
	InsertBy     string `json:"insert_by"`
	UpdateBy     string `json:"update_by"`
}
