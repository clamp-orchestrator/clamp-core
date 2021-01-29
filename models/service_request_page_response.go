package models

type ServiceRequestPageResponse struct {
	ServiceRequests []ServiceRequest `json:"serviceRequests"`
	PageNumber      int              `json:"pageNumber"`
	PageSize        int              `json:"pageSize"`
}
