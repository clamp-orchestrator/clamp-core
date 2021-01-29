package models

import (
	"time"

	"github.com/google/uuid"
)

// ServiceRequest is a structure to store the service request details
type ServiceRequest struct {
	ID           uuid.UUID              `json:"id"`
	WorkflowName string                 `json:"workflow_name"`
	Status       Status                 `json:"status"`
	CreatedAt    time.Time              `json:"created_at"`
	Payload      map[string]interface{} `json:"payload"`
	//TODO: rename to last step id executed
	CurrentStepID  int `json:"current_step_id" binding:"omitempty"`
	RequestContext RequestContext
	RequestHeaders string
}

// NewServiceRequest returns a new service request with given info
func NewServiceRequest(workflowName string, payload map[string]interface{}) ServiceRequest {
	currentTime := time.Now()
	return ServiceRequest{ID: uuid.New(), WorkflowName: workflowName, Status: STATUS_NEW, CreatedAt: currentTime, Payload: payload}
}

// A PGServiceRequest represents service request that can be persisted through go-pg
type PGServiceRequest struct {
	tableName    struct{} `pg:"service_requests"` //nolint:structcheck,unused
	ID           uuid.UUID
	WorkflowName string
	Status       Status
	CreatedAt    time.Time
	Payload      map[string]interface{} `json:"payload"`
}

// ToPgServiceRequest returns a PGServiceRequest constructed from calling ServiceRequest
func (serviceReq ServiceRequest) ToPgServiceRequest() PGServiceRequest {
	return PGServiceRequest{
		ID:           serviceReq.ID,
		WorkflowName: serviceReq.WorkflowName,
		Status:       serviceReq.Status,
		CreatedAt:    serviceReq.CreatedAt,
		Payload:      serviceReq.Payload,
	}
}

// ToServiceRequest returns a ServiceRequest constructed from calling PGServiceRequest
func (pgServReq PGServiceRequest) ToServiceRequest() ServiceRequest {
	return ServiceRequest{
		ID:           pgServReq.ID,
		WorkflowName: pgServReq.WorkflowName,
		Status:       pgServReq.Status,
		CreatedAt:    pgServReq.CreatedAt,
		Payload:      pgServReq.Payload,
	}
}
