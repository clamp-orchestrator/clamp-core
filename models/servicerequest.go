package models

import (
	"github.com/google/uuid"
)

//ServiceRequest is a structure to store the service request details
type ServiceRequest struct {
	ID           uuid.UUID `json:"id"`
	WorkflowName string    `json:"workflowName"`
	Status       Status    `json:"status"`
}

type Status string

const (
	STATUS_NEW       Status = "NEW"
	STATUS_STARTED   Status = "STARTED"
	STATUS_COMPLETED Status = "COMPLETED"
	STATUS_FAILED    Status = "FAILED"
)

func NewServiceRequest(workflowName string) ServiceRequest {
	return ServiceRequest{ID: uuid.New(), WorkflowName: workflowName, Status: STATUS_NEW}
}

type PGServiceRequest struct {
	tableName    struct{} `pg:"service_requests"`
	ID           uuid.UUID
	WorkflowName string
	Status       Status
}

func (serviceReq ServiceRequest) ToPgServiceRequest() PGServiceRequest {
	return PGServiceRequest{
		ID:           serviceReq.ID,
		WorkflowName: serviceReq.WorkflowName,
		Status:       serviceReq.Status,
	}
}

func (pgServReq PGServiceRequest) toServiceRequest() ServiceRequest {
	return ServiceRequest{
		ID:           pgServReq.ID,
		WorkflowName: pgServReq.WorkflowName,
		Status:       pgServReq.Status,
	}
}
