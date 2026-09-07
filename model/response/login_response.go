package response

type LoginResponse struct {
	IdCustomer int64  `json:"id_customer"`
	CustomerNo string `json:"customer_no"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Token      string `json:"token"`
}

type TUserModuleResponse struct {
	Module       string `json:"module"`
	CreateAccess bool   `json:"create_access"`
	ReadAccess   bool   `json:"read_access"`
	UpdateAccess bool   `json:"update_access"`
	DeleteAccess bool   `json:"delete_access"`
	ExportAccess bool   `json:"export_access"`
}

type SessionDataResponse struct {
	IdCustomer   int64                 `json:"id_customer"`
	CustomerNo   string                `json:"customer_no"`
	IdUser       int64                 `json:"id_user"`
	Username     string                `json:"username"`
	Name         string                `json:"name"`
	PhoneNo      string                `json:"phone_no"`
	EmailAddress string                `json:"email_address"`
	BloodType    string                `json:"blood_type"`
	DataModules  []TUserModuleResponse `json:"data_modules" gorm:"-"`
}

type SessionDataMainResponse struct {
	StatusResponse StatusResponse        `json:"status_response"`
	TotalPage      int64                 `json:"total_page"`
	TotalData      int64                 `json:"total_data"`
	DataHeader     interface{}           `json:"data_header"`
	Data           []SessionDataResponse `json:"data"`
}

type LoginMainResponse struct {
	StatusResponse StatusResponse  `json:"status_response"`
	TotalPage      int64           `json:"total_page"`
	TotalData      int64           `json:"total_data"`
	DataHeader     interface{}     `json:"data_header"`
	Data           []LoginResponse `json:"data"`
}
