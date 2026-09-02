package request

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Menu     string `json:"menu" binding:"required"`
	Site     string `json:"site"`
}
