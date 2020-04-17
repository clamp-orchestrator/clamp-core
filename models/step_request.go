package models

import "github.com/google/uuid"

type StepRequest struct {
	ServiceRequestId uuid.UUID              `json:"serviceRequestId"`
	StepId           int                    `json:"stepId"`
	Payload          map[string]interface{} `json:"payload"`
}
