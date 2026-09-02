package response

type TUserModuleGetResponse struct {
	Id           int64  `json:"id"`
	IdCustomer   int64  `json:"id_customer"`
	CustomerNo   string `json:"customer_no"`
	IdUser       int64  `json:"id_user"`
	Username     string `json:"username"`
	IdModule     int64  `json:"id_module"`
	Module       string `json:"module"`
	CreateAccess bool   `json:"create_access"`
	ReadAccess   bool   `json:"read_access"`
	UpdateAccess bool   `json:"update_access"`
	DeleteAccess bool   `json:"delete_access"`
	ExportAccess bool   `json:"export_access"`
	StatusData   string `json:"status_data"`
	InsertBy     string `json:"insert_by"`
	InsertTime   string `json:"insert_time"`
	UpdateBy     string `json:"update_by"`
	UpdateTime   string `json:"update_time"`
}

type TUserModuleMainResponse struct {
	StatusResponse StatusResponse           `json:"status_response"`
	TotalPage      int64                    `json:"total_page"`
	TotalData      int64                    `json:"total_data"`
	DataHeader     interface{}              `json:"data_header"`
	Data           []TUserModuleGetResponse `json:"data"`
}
