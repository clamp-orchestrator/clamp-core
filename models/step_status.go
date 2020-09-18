package models

import (
	"github.com/google/uuid"
	"time"
)

type Payload struct {
	Request  map[string]interface{} `json:"request"`
	Response map[string]interface{} `json:"response"`
}

// StepsStatus is a structure to store the service request steps details
//TODO: remove step_name field
type StepsStatus struct {
	ID               string    `json:"id"`
	ServiceRequestID uuid.UUID `json:"service_request_id"`
	WorkflowName     string    `json:"workflow_name"`
	Status           Status    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	TotalTimeInMs    int64     `json:"total_time_in_ms"`
	StepName         string    `json:"step_name"`
	Reason           string    `json:"reason"`
	Payload          Payload   `json:"payload"`
	StepID           int       `json:"step_id"`
}

func NewStepsStatus(stepStatus StepsStatus) StepsStatus {
	return StepsStatus{ID: stepStatus.ID, ServiceRequestID: stepStatus.ServiceRequestID, WorkflowName: stepStatus.WorkflowName,
		Status: STATUS_STARTED, CreatedAt: time.Now(), TotalTimeInMs: stepStatus.TotalTimeInMs, StepName: stepStatus.StepName, Reason: stepStatus.Reason}
}

// CreateStepsStatus : Entry for a given service request id and return step status details
func CreateStepsStatus(stepStatus StepsStatus) StepsStatus {
	return NewStepsStatus(stepStatus)
}

type PGStepStatus struct {
	tableName        struct{} `pg:"steps_status"`
	ID               string
	ServiceRequestID uuid.UUID
	WorkflowName     string
	Status           Status
	CreatedAt        time.Time
	TotalTimeInMs    int64
	StepName         string
	Reason           string
	Payload          Payload
	StepID           int
}

func (stepStatus StepsStatus) ToPgStepStatus() PGStepStatus {
	return PGStepStatus{
		ID:               stepStatus.ID,
		ServiceRequestID: stepStatus.ServiceRequestID,
		WorkflowName:     stepStatus.WorkflowName,
		Status:           stepStatus.Status,
		CreatedAt:        stepStatus.CreatedAt,
		TotalTimeInMs:    stepStatus.TotalTimeInMs,
		StepName:         stepStatus.StepName,
		Reason:           stepStatus.Reason,
		Payload:          stepStatus.Payload,
		StepID:           stepStatus.StepID,
	}
}

func (pgStepStatus PGStepStatus) ToStepStatus() StepsStatus {
	return StepsStatus{
		ID:               pgStepStatus.ID,
		ServiceRequestID: pgStepStatus.ServiceRequestID,
		WorkflowName:     pgStepStatus.WorkflowName,
		Status:           pgStepStatus.Status,
		CreatedAt:        pgStepStatus.CreatedAt,
		TotalTimeInMs:    pgStepStatus.TotalTimeInMs,
		StepName:         pgStepStatus.StepName,
		Reason:           pgStepStatus.Reason,
		Payload:          pgStepStatus.Payload,
		StepID:           pgStepStatus.StepID,
	}
}
