package models

// A ServiceRequestPageResponse represents service request page response
type ServiceRequestPageResponse struct {
	ServiceRequests []ServiceRequest `json:"serviceRequests"`
	PageNumber      int              `json:"pageNumber"`
	PageSize        int              `json:"pageSize"`
}
