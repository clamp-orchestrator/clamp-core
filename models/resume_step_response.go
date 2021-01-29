package models

import (
	"github.com/google/uuid"
)

// An AsyncStepResponse represents asynchronouse step's response
type AsyncStepResponse struct {
	ServiceRequestID   uuid.UUID              `json:"serviceRequestId"`
	StepID             int                    `json:"stepId"`
	Response           map[string]interface{} `json:"response"`
	Error              ClampErrorResponse     `json:"error"`
	stepStatusRecorded bool
	RequestHeaders     string
}

// SetStepStatusRecorded sets whether status is recorded or not
func (res *AsyncStepResponse) SetStepStatusRecorded(stepStatusRecorded bool) {
	res.stepStatusRecorded = stepStatusRecorded
}

// IsStepStatusRecorded returns true if the status is recorded
func (res *AsyncStepResponse) IsStepStatusRecorded() bool {
	return res.stepStatusRecorded
}

// StepProcessed to be true if step response was received internally, false if response received from external
