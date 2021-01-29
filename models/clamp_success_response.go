package models

type ClampSuccessResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func CreateSuccessResponse(code int, message string) *ClampErrorResponse {
	return &ClampErrorResponse{
		Code:    code,
		Message: message,
	}
}
