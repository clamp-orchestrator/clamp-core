package models

type ClampErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (res *ClampErrorResponse) IsNil() bool {
	return res.Code == 0 && res.Message == ""
}

func CreateErrorResponse(code int, message string) *ClampErrorResponse {
	return &ClampErrorResponse{
		Code:    code,
		Message: message,
	}
}

func EmptyErrorResponse() ClampErrorResponse {
	return ClampErrorResponse{
		Code:    0,
		Message: "",
	}
}
