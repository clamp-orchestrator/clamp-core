package models

import "github.com/google/uuid"

type RequestResponse struct {
	Request map[string]interface{}
	Response map[string]interface{}
}

type RequestContext struct {
	ServiceRequestId uuid.UUID              `json:"service_request_id"`
	WorkflowName     string                 `json:"workflow_name"`
	Payload          map[string]RequestResponse `json:"payload"`
}