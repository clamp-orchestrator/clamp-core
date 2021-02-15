package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestShouldCreateANewStepStatus(t *testing.T) {
	stepStatusRequest := &StepsStatus{
		ID:               "1",
		ServiceRequestID: uuid.New(),
		WorkflowName:     "testWF",
		Status:           StatusStarted,
		CreatedAt:        time.Now(),
		TotalTimeInMs:    0,
		StepName:         "firstStep",
		Reason:           "Success",
	}

	stepStatusResponse := CreateStepsStatus(stepStatusRequest)

	assert.NotEmpty(t, stepStatusResponse.ID)
	assert.NotEmpty(t, stepStatusResponse.ServiceRequestID)
	assert.NotNil(t, stepStatusResponse.CreatedAt)
	assert.Equal(t, stepStatusResponse.WorkflowName, stepStatusRequest.WorkflowName, fmt.Sprintf("Expected Step status name to be %s but was %s", stepStatusRequest.WorkflowName, stepStatusRequest.WorkflowName))
	assert.Equal(t, stepStatusResponse.Status, stepStatusRequest.Status, fmt.Sprintf("Expected Step status's status to be %s but was %s", stepStatusRequest.Status, stepStatusRequest.Status))
	assert.Equal(t, stepStatusResponse.StepName, stepStatusRequest.StepName, fmt.Sprintf("Expected Step status status to be %s but was %s", stepStatusRequest.StepName, stepStatusRequest.StepName))
}
