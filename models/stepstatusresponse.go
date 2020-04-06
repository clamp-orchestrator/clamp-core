package models

import (
	"github.com/google/uuid"
)

//Step Status is a structure to store the service request steps details
type StepsStatusResponse struct {
	ServiceRequestId           uuid.UUID `json:"serviceRequestId"`
	WorkflowName 			   string    `json:"workflowName"`
	Status       			   Status    `json:"status"`
	TotalTime    			   int   `json:"totalTime"`
	Steps       			   []StepResponse `json:"steps"`
	Reason   				   string    `json:"reason"`
}

type StepResponse struct {
	Name      string `json:"name"`
	Status    Status    `json:"status"`
	TimeTaken int `json:"timeTaken"`
}