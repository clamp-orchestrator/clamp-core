package models

import (
	"github.com/google/uuid"
)

type AsyncResumeStepExecutionRequest struct {
	ServiceRequestId uuid.UUID              `json:"id"`
	StepId           string                 `json:"stepId"`
	Payload          map[string]interface{} `json:"payload"`
	StepProcessed    bool                   `json:"stepProcessed" binding:default:"false"`
}