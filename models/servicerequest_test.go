package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldCreateANewServiceRequestWithUUID(t *testing.T) {
	expectedLen := 16
	servRequest := NewServiceRequest("CreateOrder")

	assert.Equal(t, expectedLen, len(servRequest.ID), fmt.Sprintf("The UUID %s should be %d chars long", servRequest.ID.String(), expectedLen))
}

func TestShouldCreateANewServiceRequestWithWorkflowName(t *testing.T) {
	expectedWorkflowName := "CreateOrder"
	servRequest := NewServiceRequest(expectedWorkflowName)

	assert.Equal(t, expectedWorkflowName, servRequest.WorkflowName, fmt.Sprintf("Expected worflow name to be %s but was %s", expectedWorkflowName, servRequest.WorkflowName))
}

func TestThatGeneratedUUIDForServiceRequestAreDifferent(t *testing.T) {
	expectedWorkflowName := "CreateOrder"
	servRequestOne := NewServiceRequest(expectedWorkflowName)
	servRequestTwo := NewServiceRequest(expectedWorkflowName)

	assert.NotEqual(t, servRequestOne, servRequestTwo, fmt.Sprintf("Expected service request UUIDs to be different but were %s and %s", servRequestOne.ID, servRequestTwo.ID))
}

func TestShouldCreateANewServiceRequestWithDefaultStatus9(t *testing.T) {
	servRequest := NewServiceRequest("CreateOrder")
	expectedStatus := STATUS_NEW
	assert.Equal(t, expectedStatus, servRequest.Status, fmt.Sprintf("Expected service request status to be equal but were %s and %s", expectedStatus, servRequest.Status))
}
