package models

//ServiceRequest is a structure to store the service request details
type WorkflowsResponse struct {
	Workflows  []Workflow `json:"workflows"`
	PageNumber int        `json:"pageNumber"`
	PageSize   int        `json:"pageSize"`
}
