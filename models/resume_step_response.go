package models

import (
	"github.com/google/uuid"
)

type AsyncStepResponse struct {
	ServiceRequestId uuid.UUID              `json:"service_request_id"`
	StepId           int                    `json:"step_id"`
	Payload          map[string]interface{} `json:"payload"`
	StepProcessed    bool                   `json:"step_processed"`
	Errors           ClampErrorResponse     `json:"errors"`
}
