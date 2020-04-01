package services

import (
	"clamp-core/models"
	"errors"
	"fmt"
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

	workflow := models.Workflow{
		Id:          "1",
		Name:        "Test",
		Description: "Test",
		Enabled:     false,
		Steps:       steps,
	}
	repo = mockGenericRepoImpl{}

	insertQueryMock = func(model interface{}) error {
		return nil
	}
	response, err := SaveWorkflow(workflow)
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.NotNil(t, response.Id)
	assert.Equal(t, workflow.Description, response.Description, fmt.Sprintf("Expected workflow description to be %s but was %s", workflow.Description, response.Description))
	assert.Equal(t, workflow.Name, response.Name, fmt.Sprintf("Expected worflow name to be %s but was %s", workflow.Name, response.Name))
	assert.Equal(t, workflow.Steps[0].Name, response.Steps[0].Name, fmt.Sprintf("Expected worflow first step name to be %s but was %s", workflow.Steps[0].Name, response.Steps[0].Name))

	insertQueryMock = func(model interface{}) error {
		return errors.New("insertion failed")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SaveWorkflow should have panicked!")
		}
	}()
	response, err = SaveWorkflow(workflow)
}
