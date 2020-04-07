package services

import (
	"clamp-core/models"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func prepareWorkflow() models.Workflow {
	steps := []models.Step{models.Step{}}

	steps[0] = models.Step{
		Id:      "firstStep",
		Name:    "firstStep",
		Enabled: true,
	}

	workflow := models.Workflow{
		Id:          "1",
		Name:        "TEST_WF",
		Description: "Test",
		Enabled:     false,
		Steps:       steps,
	}
	return workflow
}
func TestSaveWorkflow(t *testing.T) {

	workflow := prepareWorkflow()
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
	response, err = SaveWorkflow(workflow)
	assert.NotNil(t, err)
}

func TestFindWorkflowByWorkflowName(t *testing.T) {
	workflow := prepareWorkflow()

	repo = mockGenericRepoImpl{}

	whereQueryMock = func(model interface{}, cond string, params ...interface{}) error {
		workFlowReq := model.(*models.Workflow)
		workFlowReq.Name = workflow.Name
		workFlowReq.Steps = workflow.Steps
		return nil
	}
	resp, err := FindWorkflowByName(workflow.Name)
	assert.Nil(t, err)
	assert.Equal(t, workflow.Name, resp.Name)
	assert.NotNil(t, resp.Steps)

	whereQueryMock = func(model interface{}, cond string, params ...interface{}) error {
		return errors.New("select query failed")
	}
	_, err = FindWorkflowByName(workflow.Name)
	assert.NotNil(t, err)
}
