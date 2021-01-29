package models

// A ClampErrorResponse represents an error returned from Clamp API response
type ClampErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// IsNil returns true if it doesn't represent any error
func (res *ClampErrorResponse) IsNil() bool {
	return res.Code == 0 && res.Message == ""
}

// CreateErrorResponse returns a ClampErrorResponse with specified code and message
func CreateErrorResponse(code int, message string) *ClampErrorResponse {
	return &ClampErrorResponse{
		Code:    code,
		Message: message,
	}
}

// EmptyErrorResponse returns nil ClampErrorResponse
func EmptyErrorResponse() ClampErrorResponse {
	return ClampErrorResponse{
		Code:    0,
		Message: "",
	}
}
