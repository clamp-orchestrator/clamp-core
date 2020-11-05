package models

type WorkflowsPageResponse struct {
	Workflows      []Workflow `json:"workflows"`
	PageNumber     int        `json:"pageNumber"`
	PageSize       int        `json:"pageSize"`
	TotalWorkflows int        `json:"totalWorkflows"`
}
