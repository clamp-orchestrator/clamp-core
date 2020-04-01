package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldCreateANewWorkflow(t *testing.T) {
	steps := []Step{Step{}}

	steps[0] = Step{
		Id:      "firstStep",
		Name:    "firstStep",
		Enabled: true,
	}
	workflow := Workflow{
		Id:          "1",
		Name:        "Test",
		Description: "Test",
		Enabled:     false,
		Steps:       steps,
	}

	serviceFlowRequest := workflow
	workflowResponse := CreateWorkflow(serviceFlowRequest)

	assert.NotEmpty(t, workflowResponse.Id)
	assert.Equal(t, serviceFlowRequest.Description, workflowResponse.Description, fmt.Sprintf("Expected workflow description to be %s but was %s", serviceFlowRequest.Description, workflowResponse.Description))
	assert.Equal(t, serviceFlowRequest.Name, workflowResponse.Name, fmt.Sprintf("Expected worflow name to be %s but was %s", serviceFlowRequest.Name, workflowResponse.Name))
	assert.Equal(t, serviceFlowRequest.Steps[0].Name, workflowResponse.Steps[0].Name, fmt.Sprintf("Expected worflow first step name to be %s but was %s", serviceFlowRequest.Steps[0].Name, workflowResponse.Steps[0].Name))
}
