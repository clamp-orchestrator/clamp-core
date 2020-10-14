package services

import (
	"clamp-core/models"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepareWorkflow() models.Workflow {
	steps := []models.Step{models.Step{}}

	steps[0] = models.Step{
		Name:    "firstStep",
		Enabled: true,
	}

	workflow := models.Workflow{
		ID:          "1",
		Name:        "TEST_WF",
		Description: "Test",
		Enabled:     false,
		Steps:       steps,
	}
	return workflow
}
func TestSaveWorkflow(t *testing.T) {

	workflow := prepareWorkflow()
	SaveWorkflowMock = func(workflow models.Workflow) (models.Workflow, error) {
		return workflow, nil
	}
	response, err := SaveWorkflow(workflow)
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.NotNil(t, response.ID)
	assert.Equal(t, workflow.Description, response.Description, fmt.Sprintf("Expected workflow description to be %s but was %s", workflow.Description, response.Description))
	assert.Equal(t, workflow.Name, response.Name, fmt.Sprintf("Expected worflow name to be %s but was %s", workflow.Name, response.Name))
	assert.Equal(t, workflow.Steps[0].Name, response.Steps[0].Name, fmt.Sprintf("Expected worflow first step name to be %s but was %s", workflow.Steps[0].Name, response.Steps[0].Name))

	SaveWorkflowMock = func(workflow models.Workflow) (models.Workflow, error) {
		return workflow, errors.New("insertion failed")
	}
	response, err = SaveWorkflow(workflow)
	assert.NotNil(t, err)
}

func TestFindWorkflowByWorkflowName(t *testing.T) {
	workflow := prepareWorkflow()

	findWorkflowByNameMock = func(workflowName string) (models.Workflow, error) {
		return workflow, nil
	}
	resp, err := FindWorkflowByName(workflow.Name)
	assert.Nil(t, err)
	assert.Equal(t, workflow.Name, resp.Name)
	assert.NotNil(t, resp.Steps)
	findWorkflowByNameMock = func(workflowName string) (models.Workflow, error) {
		return workflow, errors.New("select query failed")
	}
	_, err = FindWorkflowByName(workflow.Name)
	assert.NotNil(t, err)
}

func TestGetWorkflowsWithoutSortByArgs(t *testing.T) {
	workflow := prepareWorkflow()
	sortBy := make(map[string]string)
	sortOrder := []string{}
	getWorkflowsMock = func(pageNumber int, pageSize int, sortBy map[string]string, sortOrder []string) ([]models.Workflow, error) {
		return []models.Workflow{workflow}, nil
	}
	resp, err := GetWorkflows(0, 1, sortBy, sortOrder)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(resp))
}

func TestGetWorkflowsWithSortByArgs(t *testing.T) {
	workflow := prepareWorkflow()
	sortBy := map[string]string{
		"id": "asc",
	}
	sortOrder := []string{"id"}
	getWorkflowsMock = func(pageNumber int, pageSize int, sortBy map[string]string, sortOrder []string) ([]models.Workflow, error) {
		return []models.Workflow{workflow}, nil
	}
	resp, err := GetWorkflows(0, 1, sortBy, sortOrder)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(resp))
}
