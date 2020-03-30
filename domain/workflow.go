package domain

import (
	"github.com/google/uuid"
)

//Workflow is a structure to store the service request details
type Workflow struct {
	ID           uuid.UUID
	ServiceFlow ServiceFlow
}
type ServiceFlow struct {
	Description string `json:"description"`
	FlowMode    string `json:"flowMode"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Steps       Steps  `json:"steps"`
}

type Steps struct {
	Step []Step `json:"step"`
}

type Step struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

//Create a new work flow for a given service flow and return service flow details
func CreateWorkflow(serviceFlowRequest Workflow) Workflow {
	return newServiceFlow(serviceFlowRequest)
}


func newServiceFlow(workflow Workflow) Workflow {
	return Workflow{ServiceFlow: workflow.ServiceFlow}
}
