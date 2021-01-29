package models

import "github.com/google/uuid"

// A StepContext holds context data related to workflow step
type StepContext struct {
	Request         map[string]interface{}
	RequestHeaders  string
	ResponseHeaders string
	Response        map[string]interface{}
	StepSkipped     bool
}

// A RequestContext holds context data related to service request
type RequestContext struct {
	ServiceRequestID uuid.UUID               `json:"service_request_id"`
	WorkflowName     string                  `json:"workflow_name"`
	StepsContext     map[string]*StepContext `json:"step_context"`
}

// SetStepRequestToContext sets given request payload to the specified step
func (ctx *RequestContext) SetStepRequestToContext(stepName string, requestPayload map[string]interface{}) {
	stepContext := ctx.StepsContext[stepName]
	stepContext.Request = requestPayload
}

// SetStepResponseToContext sets given response payload to the specified step
func (ctx *RequestContext) SetStepResponseToContext(stepName string, responsePayload map[string]interface{}) {
	stepContext := ctx.StepsContext[stepName]
	stepContext.Response = responsePayload
}

// SetStepRequestHeadersToContext sets given request headers to the specified step
func (ctx *RequestContext) SetStepRequestHeadersToContext(stepName string, requestHeaders string) {
	stepContext := ctx.StepsContext[stepName]
	stepContext.RequestHeaders = requestHeaders
}

// SetStepResponseHeadersToContext sets given response headers to the specified step
func (ctx *RequestContext) SetStepResponseHeadersToContext(stepName string, responseHeaders string) {
	stepContext := ctx.StepsContext[stepName]
	stepContext.ResponseHeaders = responseHeaders
}

// GetStepRequestHeadersFromContext returns request header of the specified step
func (ctx RequestContext) GetStepRequestHeadersFromContext(stepName string) string {
	stepContext := ctx.StepsContext[stepName]
	return stepContext.RequestHeaders
}

// GetStepResponseHeadersFromContext returns response header of the specified step
func (ctx RequestContext) GetStepResponseHeadersFromContext(stepName string) string {
	stepContext := ctx.StepsContext[stepName]
	return stepContext.ResponseHeaders
}

// GetStepRequestFromContext returns request of the specified step
func (ctx RequestContext) GetStepRequestFromContext(stepName string) map[string]interface{} {
	stepContext := ctx.StepsContext[stepName]
	return stepContext.Request
}

// GetStepResponseFromContext returns response of the specified step
func (ctx RequestContext) GetStepResponseFromContext(stepName string) map[string]interface{} {
	stepContext := ctx.StepsContext[stepName]
	return stepContext.Response
}
