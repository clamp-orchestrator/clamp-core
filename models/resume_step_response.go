package models

import (
	"github.com/google/uuid"
)

type AsyncStepResponse struct {
	ServiceRequestId uuid.UUID              `json:"serviceRequestId"`
	StepId           int                    `json:"stepId"`
	Payload          map[string]interface{} `json:"payload"`
	StepProcessed    bool                   `json:"stepProcessed" binding:default:"false"`
	Errors           ClampErrorResponse     `json:"errors"`
}
