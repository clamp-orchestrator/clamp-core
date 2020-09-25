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
		ServiceRequestID: uuid.UUID{},
		WorkflowName:     "test_wf",
		StepsContext:     map[string]*StepContext{"step1": {nil, "Content-Type:application/json;","",nil, false}},
	}
	context.SetStepRequestToContext("step1", payload)
	assert.NotNil(t, context.StepsContext["step1"].Request)
	assert.Equal(t, payload, context.StepsContext["step1"].Request)
}

func TestShouldGetStepRequestFromContext(t *testing.T) {
	payload := map[string]interface{}{"k": "v"}
	context := RequestContext{
		ServiceRequestID: uuid.UUID{},
		WorkflowName:     "test_wf",
		StepsContext:     map[string]*StepContext{"step1": {nil, "","",nil, false}},
	}
	context.SetStepRequestToContext("step1", payload)
	assert.NotNil(t, context.StepsContext["step1"].Request)
	assert.Equal(t, payload, context.GetStepRequestFromContext("step1"))
}

func TestShouldSetStepResponseToContext(t *testing.T) {
	payload := map[string]interface{}{"k": "v"}
	context := RequestContext{
		ServiceRequestID: uuid.UUID{},
		WorkflowName:     "test_wf",
		StepsContext:     map[string]*StepContext{"step1": {nil, "","",nil, false}},
	}
	context.SetStepResponseToContext("step1", payload)
	assert.NotNil(t, context.StepsContext["step1"].Response)
	assert.Equal(t, payload, context.StepsContext["step1"].Response)
}

func TestShouldGetStepResponseFromContext(t *testing.T) {
	payload := map[string]interface{}{"k": "v"}
	context := RequestContext{
		ServiceRequestID: uuid.UUID{},
		WorkflowName:     "test_wf",
		StepsContext:     map[string]*StepContext{"step1": {nil, "","",nil, false}},
	}
	context.SetStepResponseToContext("step1", payload)
	assert.NotNil(t, context.StepsContext["step1"].Response)
	assert.Equal(t, payload, context.GetStepResponseFromContext("step1"))
}

func TestShouldSetStepRequestHeadersToContext(t *testing.T) {
	requestHeaders := "Content-Type:application/json"
	context := RequestContext{
		ServiceRequestID: uuid.UUID{},
		WorkflowName:     "test_wf",
		StepsContext:     map[string]*StepContext{"step1": {nil, "Content-Type:application/json;","",nil, false}},
	}
	context.SetStepRequestHeadersToContext("step1", requestHeaders)
	assert.NotNil(t, context.StepsContext["step1"].RequestHeaders)
	assert.Equal(t, requestHeaders, context.StepsContext["step1"].RequestHeaders)
}

func TestShouldGetStepRequestHeadersToContext(t *testing.T) {
	requestHeaders := "Content-Type:application/json"
	context := RequestContext{
		ServiceRequestID: uuid.UUID{},
		WorkflowName:     "test_wf",
		StepsContext:     map[string]*StepContext{"step1": {nil, "Content-Type:application/json;","",nil, false}},
	}
	context.SetStepRequestHeadersToContext("step1", requestHeaders)
	assert.NotNil(t, context.StepsContext["step1"].RequestHeaders)
	assert.Equal(t, requestHeaders, context.GetStepRequestHeadersFromContext("step1"))
}

func TestShouldSetStepResponseHeadersToContext(t *testing.T) {
	responseHeaders := "Content-Type:application/json"
	context := RequestContext{
		ServiceRequestID: uuid.UUID{},
		WorkflowName:     "test_wf",
		StepsContext:     map[string]*StepContext{"step1": {nil, "Content-Type:application/json;","Content-Type:application/json;",nil, false}},
	}
	context.SetStepResponseHeadersToContext("step1", responseHeaders)
	assert.NotNil(t, context.StepsContext["step1"].ResponseHeaders)
	assert.Equal(t, responseHeaders, context.StepsContext["step1"].ResponseHeaders)
}

func TestShouldGetStepResponseHeadersToContext(t *testing.T) {
	responseHeaders := "Content-Type:application/json;"
	context := RequestContext{
		ServiceRequestID: uuid.UUID{},
		WorkflowName:     "test_wf",
		StepsContext:     map[string]*StepContext{"step1": {nil, "Content-Type:application/json;","Content-Type:application/json;",nil, false}},
	}
	context.SetStepRequestHeadersToContext("step1", responseHeaders)
	assert.NotNil(t, context.StepsContext["step1"].ResponseHeaders)
	assert.Equal(t, responseHeaders, context.GetStepResponseHeadersFromContext("step1"))
}

