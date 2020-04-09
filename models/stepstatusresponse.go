package models

import (
	"github.com/google/uuid"
)

//Step Status Response is a structure to display status of Service request id to the users
type StepsStatusResponse struct {
	ServiceRequestId           uuid.UUID `json:"serviceRequestId"`
	WorkflowName 			   string    `json:"workflowName"`
	Status       			   Status    `json:"status"`
	TotalTimeInMs    		   int64   `json:"totalTimeInMs"`
	Steps       			   []StepResponse `json:"steps"`
	Reason   				   string    `json:"reason"`
}

type StepResponse struct {
	Name      string `json:"name"`
	Status    Status    `json:"status"`
	TimeTaken int64 `json:"timeTaken"`
	Payload      map[string]interface{}   `json:"payload"`
}