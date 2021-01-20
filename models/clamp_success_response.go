package models

// A ClampSuccessResponse represents a successful response from Clamp API
type ClampSuccessResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// CreateSuccessResponse returns a ClampErrorResponse with specified code and message
func CreateSuccessResponse(code int, message string) *ClampErrorResponse {
	return &ClampErrorResponse{
		Code:    code,
		Message: message,
	}
}
