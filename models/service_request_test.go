package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepareRequestPayload() map[string]interface{} {
	var payload = map[string]interface{}{
		"id1": "val1",
		"id2": "val2",
	}
	return payload
}
func TestShouldCreateANewServiceRequestWithUUID(t *testing.T) {
	expectedLen := 16
	payload := prepareRequestPayload()

	servRequest := NewServiceRequest("CreateOrder", payload)

	assert.Equal(t, expectedLen, len(servRequest.ID), fmt.Sprintf("The UUID %s should be %d chars long", servRequest.ID.String(), expectedLen))
}

func TestShouldCreateANewServiceRequestWithWorkflowName(t *testing.T) {
	expectedWorkflowName := "CreateOrder"

	payload := prepareRequestPayload()
	servRequest := NewServiceRequest(expectedWorkflowName, payload)

	assert.Equal(t, expectedWorkflowName, servRequest.WorkflowName, fmt.Sprintf("Expected worflow name to be %s but was %s", expectedWorkflowName, servRequest.WorkflowName))
}

func TestThatGeneratedUUIDForServiceRequestAreDifferent(t *testing.T) {
	expectedWorkflowName := "CreateOrder"

	payload := prepareRequestPayload()
	servRequestOne := NewServiceRequest(expectedWorkflowName, payload)
	servRequestTwo := NewServiceRequest(expectedWorkflowName, payload)

	assert.NotEqual(t, servRequestOne, servRequestTwo, fmt.Sprintf("Expected service request UUIDs to be different but were %s and %s", servRequestOne.ID, servRequestTwo.ID))
}

func TestShouldCreateANewServiceRequestWithDefaultStatus9(t *testing.T) {

	payload := prepareRequestPayload()
	servRequest := NewServiceRequest("CreateOrder", payload)

	expectedStatus := StatusNew
	assert.Equal(t, expectedStatus, servRequest.Status, fmt.Sprintf("Expected service request status to be equal but were %s and %s", expectedStatus, servRequest.Status))
}
