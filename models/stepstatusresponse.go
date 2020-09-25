package models

import (
	"github.com/google/uuid"
)

//Step Status Response is a structure to display status of Service request id to the users
type ServiceRequestStatusResponse struct {
	ServiceRequestID uuid.UUID            `json:"service_request_id"`
	WorkflowName     string               `json:"workflow_name"`
	Status           Status               `json:"status"`
	TotalTimeInMs    int64                `json:"total_time_in_ms"`
	Steps            []StepStatusResponse `json:"steps"`
	Reason           string               `json:"reason"`
}

type StepStatusResponse struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Status    Status  `json:"status"`
	TimeTaken int64   `json:"time_taken"`
	Payload   Payload `json:"payload"`
}
