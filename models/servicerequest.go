package models

import (
	"github.com/google/uuid"
	"time"
)

//ServiceRequest is a structure to store the service request details
type ServiceRequest struct {
	ID           uuid.UUID `json:"id"`
	WorkflowName string    `json:"workflowName"`
	Status       Status    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Status string

const (
	STATUS_NEW       Status = "NEW"
	STATUS_STARTED   Status = "STARTED"
	STATUS_COMPLETED Status = "COMPLETED"
	STATUS_FAILED    Status = "FAILED"
)

func NewServiceRequest(workflowName string) ServiceRequest {
	currentTime := time.Now()
	return ServiceRequest{ID: uuid.New(), WorkflowName: workflowName, Status: STATUS_NEW, CreatedAt: currentTime}
}

type PGServiceRequest struct {
	tableName    struct{} `pg:"service_requests"`
	ID           uuid.UUID
	WorkflowName string
	Status       Status
	CreatedAt    time.Time
}

func (serviceReq ServiceRequest) ToPgServiceRequest() PGServiceRequest {
	return PGServiceRequest{
		ID:           serviceReq.ID,
		WorkflowName: serviceReq.WorkflowName,
		Status:       serviceReq.Status,
		CreatedAt:    serviceReq.CreatedAt,
	}
}

func (pgServReq PGServiceRequest) ToServiceRequest() ServiceRequest {
	return ServiceRequest{
		ID:           pgServReq.ID,
		WorkflowName: pgServReq.WorkflowName,
		Status:       pgServReq.Status,
		CreatedAt:    pgServReq.CreatedAt,
	}
}
