package request

type OnboardRequest struct {
	Serial  string `json:"serial" binding:"required"`
	License string `json:"license" binding:"required"`
	Menu    string `json:"menu" binding:"required"`
}
