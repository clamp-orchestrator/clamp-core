package models

import "github.com/google/uuid"

type AsyncStepRequest struct {
	Step             Step                   `json:"step"`
	Payload          map[string]interface{} `json:"payload"`
	ServiceRequestId uuid.UUID              `json:"service_request_id"`
	WorkflowName     string                 `json:"workflow_name"`
}
