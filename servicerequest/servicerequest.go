package servicerequest

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

//Create a new service request for a given flow name and return service request details
func Create(workflowName string) ServiceRequest {
	return newServiceRequest(workflowName)
}

func newServiceRequest(workflowName string) ServiceRequest {
	return ServiceRequest{ID: uuid.New(), WorkflowName: workflowName, Status: STATUS_NEW}
}
