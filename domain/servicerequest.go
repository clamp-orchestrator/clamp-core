package domain

import (
	"github.com/google/uuid"
)

//ServiceRequest is a structure to store the service request details
type ServiceRequest struct {
	ID           uuid.UUID `json:"id"`
	WorkflowName string    `json:"workflowName"`
}

//Create a new service request for a given flow name and return service request details
func Create(workflowName string) ServiceRequest {
	return newServiceRequest(workflowName)
}

func newServiceRequest(workflowName string) ServiceRequest {
	return ServiceRequest{ID: uuid.New(), WorkflowName: workflowName}
}
