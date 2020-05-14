package models

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	var employee = make(map[string]StepContext)
	request := StepContext{
		Request:  prepareRequestPayload(),
		Response: prepareRequestPayload(),
	}
	employee["Mark"] = request
	employee["Sandy"] = request
	fmt.Println(employee)
}

func TestShouldSetStepRequestToContext(t *testing.T) {
	payload := map[string]interface{}{"k": "v"}
	context := RequestContext{
		ServiceRequestId: uuid.UUID{},
		WorkflowName:     "test_wf",
		StepsContext:     map[string]*StepContext{"step1": {nil, nil, false}},
	}
	context.SetStepRequestToContext("step1", payload)
	assert.NotNil(t, context.StepsContext["step1"].Request)
	assert.Equal(t, payload, context.StepsContext["step1"].Request)
}

func TestShouldGetStepRequestFromContext(t *testing.T) {
	payload := map[string]interface{}{"k": "v"}
	context := RequestContext{
		ServiceRequestId: uuid.UUID{},
		WorkflowName:     "test_wf",
		StepsContext:     map[string]*StepContext{"step1": {nil, nil, false}},
	}
	context.SetStepRequestToContext("step1", payload)
	assert.NotNil(t, context.StepsContext["step1"].Request)
	assert.Equal(t, payload, context.GetStepRequestFromContext("step1"))
}

func TestShouldSetStepResponseToContext(t *testing.T) {
	payload := map[string]interface{}{"k": "v"}
	context := RequestContext{
		ServiceRequestId: uuid.UUID{},
		WorkflowName:     "test_wf",
		StepsContext:     map[string]*StepContext{"step1": {nil, nil, false}},
	}
	context.SetStepResponseToContext("step1", payload)
	assert.NotNil(t, context.StepsContext["step1"].Response)
	assert.Equal(t, payload, context.StepsContext["step1"].Response)
}

func TestShouldGetStepResponseFromContext(t *testing.T) {
	payload := map[string]interface{}{"k": "v"}
	context := RequestContext{
		ServiceRequestId: uuid.UUID{},
		WorkflowName:     "test_wf",
		StepsContext:     map[string]*StepContext{"step1": {nil, nil, false}},
	}
	context.SetStepResponseToContext("step1", payload)
	assert.NotNil(t, context.StepsContext["step1"].Response)
	assert.Equal(t, payload, context.GetStepResponseFromContext("step1"))
}
