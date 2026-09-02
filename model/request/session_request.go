package request

type SessionRequest struct {
	Username   string `json:"username" binding:"required"`
	CustomerID int64  `json:"customer_id" binding:"required"`
	CustomerNo string `json:"customer_no" binding:"required"`
	Site       string `json:"site"`
	Menu       string `json:"menu" binding:"required"`
}

type SessionParamRequest struct {
	Site string `form:"site"`
}

type RedisRequest struct {
	Username    string `json:"username" binding:"required"`
	CustomerID  int64  `json:"customer_id" binding:"required"`
	CustomerNo  string `json:"customer_no" binding:"required"`
	Site        string `json:"site" binding:"required"`
	Menu        string `json:"menu" binding:"required"`
	CompanyCode string `json:"company_code" binding:"required"`
	CompanyName string `json:"company_name" binding:"required"`
	Permission  string `json:"permission"`
}

type GoogleLoginRequest struct {
	IdToken string `json:"id_token" binding:"required"`
}

type FacebookLoginRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}
