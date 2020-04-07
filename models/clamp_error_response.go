package models

type ClampErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func CreateErrorResponse(code int, message string) *ClampErrorResponse {
	return &ClampErrorResponse{
		Code:    code,
		Message: message,
	}
}
