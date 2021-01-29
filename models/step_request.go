package models

import "github.com/google/uuid"

// A StepRequest represents a request initiated from workflow step
type StepRequest struct {
	ServiceRequestID uuid.UUID              `json:"service_request_id"`
	StepID           int                    `json:"step_id"`
	Payload          map[string]interface{} `json:"payload"`
	Headers          string
}

// NewStepRequest returns a new workflow step request with the given info
func NewStepRequest(serviceRequestID uuid.UUID, stepID int, payload map[string]interface{}, headers string) *StepRequest {
	return &StepRequest{ServiceRequestID: serviceRequestID, StepID: stepID, Payload: payload, Headers: headers}
}
