package services

import (
	"clamp-core/models"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSaveWorkflow(t *testing.T) {
	steps := []models.Step{models.Step{}}

	steps[0] = models.Step{
		Id:      "firstStep",
		Name:    "firstStep",
		Enabled: true,
	}

	serviceFlowRequest := models.Workflow{
		ID:          uuid.New(),
		ServiceFlow: models.ServiceFlow{
			Description: "Test",
			FlowMode:    "None",
			Id:          "1",
			Name:        "Test",
			Enabled:     true,
			Steps:       models.Steps{
				Step: steps,
			},
		},
	}
	repo = mockGenericRepoImpl{}

	insertQueryMock = func(model interface{}) error {
		return nil
	}
	response, err := SaveServiceFlow(serviceFlowRequest)
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.NotNil(t, response.ID)
	assert.Equal(t, serviceFlowRequest.ServiceFlow.Description, response.ServiceFlow.Description, fmt.Sprintf("Expected workflow description to be %s but was %s", serviceFlowRequest.ServiceFlow.Description, response.ServiceFlow.Description))
	assert.Equal(t, serviceFlowRequest.ServiceFlow.Name, response.ServiceFlow.Name, fmt.Sprintf("Expected worflow name to be %s but was %s", serviceFlowRequest.ServiceFlow.Name, response.ServiceFlow.Name))
	assert.Equal(t, serviceFlowRequest.ServiceFlow.Steps.Step[0].Name, response.ServiceFlow.Steps.Step[0].Name, fmt.Sprintf("Expected worflow first step name to be %s but was %s", serviceFlowRequest.ServiceFlow.Steps.Step[0].Name, response.ServiceFlow.Steps.Step[0].Name))

	insertQueryMock = func(model interface{}) error {
		return errors.New("insertion failed")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SaveWorkflow should have panicked!")
		}
	}()
	response, err = SaveServiceFlow(serviceFlowRequest)
}
