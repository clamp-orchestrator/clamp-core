package models

import "github.com/google/uuid"

type StepContext struct {
	Request  map[string]interface{}
	Response map[string]interface{}
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

func (ctx RequestContext) GetStepRequestFromContext(stepName string) map[string]interface{} {
	stepContext := ctx.StepsContext[stepName]
	return stepContext.Request
}
func (ctx RequestContext) GetStepResponseFromContext(stepName string) map[string]interface{} {
	stepContext := ctx.StepsContext[stepName]
	return stepContext.Response
}
