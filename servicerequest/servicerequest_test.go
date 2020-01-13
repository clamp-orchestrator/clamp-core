package servicerequest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldCreateANewServiceRequestWithUUID(t *testing.T) {
	expectedLen := 36
	servRequest := Create("CreateOrder")
	strUUID := servRequest.ID.String()

	assert.Equal(t, expectedLen, len(strUUID), fmt.Sprintf("The UUID %s should be %d chars long", strUUID, expectedLen))
}

func TestShouldCreateANewServiceRequestWithWorkflowName(t *testing.T) {
	expectedWorkflowName := "CreateOrder"
	servRequest := Create(expectedWorkflowName)

	assert.Equal(t, expectedWorkflowName, servRequest.workflowName, fmt.Sprintf("Expected worflow name to be %s but was %s", expectedWorkflowName, servRequest.workflowName))
}

func TestThatGeneratedUUIDForServiceRequestAreDifferent(t *testing.T) {
	expectedWorkflowName := "CreateOrder"
	servRequestOne := Create(expectedWorkflowName)
	servRequestTwo := Create(expectedWorkflowName)

	assert.NotEqual(t, servRequestOne, servRequestTwo, fmt.Sprintf("Expected service request UUIDs to be different but were %s and %s", servRequestOne.ID, servRequestTwo.ID))
}
