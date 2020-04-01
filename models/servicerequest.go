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
	StartTime    time.Time    `json:"startTime"`
	EndTime      time.Time   `json:"endTime"`
	TotalTimeElapsedMs      int    `json:"totalTimeElapsedMs"`
	Steps      Step    `json:"steps"`
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
	return ServiceRequest{ID: uuid.New(), WorkflowName: workflowName, Status: STATUS_NEW, StartTime:currentTime, EndTime:currentTime, TotalTimeElapsedMs:0, Steps: Step{}}
}

type PGServiceRequest struct {
	tableName    struct{} `pg:"service_requests"`
	ID           uuid.UUID
	WorkflowName string
	Status       Status
	StartTime    time.Time
	EndTime      time.Time
	TotalTimeElapsedMs      int
	Steps      	 Step
}

func (serviceReq ServiceRequest) ToPgServiceRequest() PGServiceRequest {
	return PGServiceRequest{
		ID:           serviceReq.ID,
		WorkflowName: serviceReq.WorkflowName,
		Status:       serviceReq.Status,
		StartTime:       serviceReq.StartTime,
		EndTime:       serviceReq.EndTime,
		TotalTimeElapsedMs:       serviceReq.TotalTimeElapsedMs,
		Steps:       serviceReq.Steps,
	}
}

func (pgServReq PGServiceRequest) toServiceRequest() ServiceRequest {
	return ServiceRequest{
		ID:           pgServReq.ID,
		WorkflowName: pgServReq.WorkflowName,
		Status:       pgServReq.Status,
		StartTime:       pgServReq.StartTime,
		EndTime:       pgServReq.EndTime,
		TotalTimeElapsedMs:       pgServReq.TotalTimeElapsedMs,
		Steps:       pgServReq.Steps,
	}
}
