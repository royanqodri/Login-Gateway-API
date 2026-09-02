package response

type TUserGetResponse struct {
	Id           int64  `json:"id"`
	IdCustomer   int64  `json:"id_customer"`
	CustomerNo   string `json:"customer_no"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	PhoneNo      string `json:"phone_no"`
	EmailAddress string `json:"email_address"`
	BloodType    string `json:"blood_type"`
	StatusData   string `json:"status_data"`
	InsertBy     string `json:"insert_by"`
	InsertTime   string `json:"insert_time"`
	UpdateBy     string `json:"update_by"`
	UpdateTime   string `json:"update_time"`
}

type TUserPostResponse struct {
	Type         string `json:"type"`
	Id           int64  `json:"id"`
	IdCustomer   int64  `json:"id_customer"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	PhoneNo      string `json:"phone_no"`
	EmailAddress string `json:"email_address"`
	BloodType    string `json:"blood_type"`
	StatusData   string `json:"status_data"`
	InsertBy     string `json:"insert_by"`
	UpdateBy     string `json:"update_by"`
	StatusCode   int    `json:"status_code"`
	StatusText   string `json:"status_text"`
	Message      string `json:"message"`
}

type TUserMainResponse struct {
	StatusResponse StatusResponse     `json:"status_response"`
	TotalPage      int64              `json:"total_page"`
	TotalData      int64              `json:"total_data"`
	DataHeader     interface{}        `json:"data_header"`
	Data           []TUserGetResponse `json:"data"`
}
