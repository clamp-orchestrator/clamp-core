package models

import (
	"fmt"
	"github.com/google/uuid"
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

	serviceFlow := ServiceFlow{
		Description: "Test",
		FlowMode:    "None",
		Id:          "1",
		Name:        "Test",
		Enabled:     true,
		Steps: Steps{
			Step: steps,
		},
	}

	serviceFlowRequest := Workflow{ID: uuid.New(), ServiceFlow: serviceFlow}
	workflowResponse := CreateWorkflow(serviceFlowRequest)

	assert.NotEmpty(t, workflowResponse.ID)
	assert.Equal(t, serviceFlowRequest.ServiceFlow.Description, workflowResponse.ServiceFlow.Description, fmt.Sprintf("Expected workflow description to be %s but was %s", serviceFlowRequest.ServiceFlow.Description, workflowResponse.ServiceFlow.Description))
	assert.Equal(t, serviceFlowRequest.ServiceFlow.Name, workflowResponse.ServiceFlow.Name, fmt.Sprintf("Expected worflow name to be %s but was %s", serviceFlowRequest.ServiceFlow.Name, workflowResponse.ServiceFlow.Name))
	assert.Equal(t, serviceFlowRequest.ServiceFlow.Steps.Step[0].Name, workflowResponse.ServiceFlow.Steps.Step[0].Name, fmt.Sprintf("Expected worflow first step name to be %s but was %s", serviceFlowRequest.ServiceFlow.Steps.Step[0].Name, workflowResponse.ServiceFlow.Steps.Step[0].Name))
}
