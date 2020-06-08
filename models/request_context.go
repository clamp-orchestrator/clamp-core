package models

import "github.com/google/uuid"

type StepContext struct {
	Request  map[string]interface{}
	RequestHeaders  string
	ResponseHeaders  string
	Response map[string]interface{}
	StepSkipped bool
}

type RequestContext struct {
	ServiceRequestId uuid.UUID               `json:"service_request_id"`
	WorkflowName     string                  `json:"workflow_name"`
	StepsContext     map[string]*StepContext `json:"step_context"`
}

func (ctx *RequestContext) SetStepRequestToContext(stepName string, requestPayload map[string]interface{}) {
	stepContext := ctx.StepsContext[stepName]
	stepContext.Request = requestPayload
}

func (ctx *RequestContext) SetStepResponseToContext(stepName string, responsePayload map[string]interface{}) {
	stepContext := ctx.StepsContext[stepName]
	stepContext.Response = responsePayload
}

func (ctx *RequestContext) SetStepRequestHeadersToContext(stepName string, requestHeaders string) {
	stepContext := ctx.StepsContext[stepName]
	stepContext.RequestHeaders = requestHeaders
}

func (ctx *RequestContext) SetStepResponseHeadersToContext(stepName string, responseHeaders string) {
	stepContext := ctx.StepsContext[stepName]
	stepContext.ResponseHeaders = responseHeaders
}

func (ctx RequestContext) GetStepRequestHeadersFromContext(stepName string) string {
	stepContext := ctx.StepsContext[stepName]
	return stepContext.RequestHeaders
}

func (ctx RequestContext) GetStepResponseHeadersFromContext(stepName string) string {
	stepContext := ctx.StepsContext[stepName]
	return stepContext.ResponseHeaders
}


func (ctx RequestContext) GetStepRequestFromContext(stepName string) map[string]interface{} {
	stepContext := ctx.StepsContext[stepName]
	return stepContext.Request
}
func (ctx RequestContext) GetStepResponseFromContext(stepName string) map[string]interface{} {
	stepContext := ctx.StepsContext[stepName]
	return stepContext.Response
}
