package models

import "github.com/google/uuid"

type StepRequest struct {
	ServiceRequestId uuid.UUID              `json:"service_request_id"`
	StepId           int                    `json:"step_id"`
	Payload          map[string]interface{} `json:"payload"`
}
