package models

import (
	"clamp-core/executors"
	"clamp-core/hooks"
	"clamp-core/transform"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// Val represents a value of a workflow step
type Val interface {
}

// RequestTransform represents the specification of request transformation
type RequestTransform interface {
}

// Step represents a workflow step
type Step struct {
	ID               int              `json:"id"`
	Name             string           `json:"name" binding:"required"`
	Type             string           `json:"type"`
	Mode             string           `json:"mode" binding:"required,oneof=HTTP AMQP KAFKA"`
	Val              Val              `json:"val" binding:"required"`
	Transform        bool             `json:"transform"`
	Enabled          bool             `json:"enabled"`
	When             string           `json:"when"`
	TransformFormat  string           `json:"transformFormat"`
	RequestTransform RequestTransform `json:"requestTransform"`
	canStepExecute   bool
	OnFailure        []Step `json:"onFailure"`
	//shouldStepExecute func(whenCondition string, stepRequest map[string]interface{}, prefix string) (canStepExecute bool, _ error)
	//transformRequest  func(stepRequest map[string]interface{}, prefix string) (map[string]interface{}, error)
}

// DidStepExecute returns true if the step can be executed
func (step *Step) DidStepExecute() bool {
	return step.canStepExecute
}

// PreStepExecution prepares the step for execution with the given info
func (step *Step) PreStepExecution(contextPayload map[string]*StepContext, prefix string) (err error) {
	step.canStepExecute = true
	stepRequestResponsePayload := make(map[string]interface{})

	if step.When != "" {
		for s, stepRequestResponse := range contextPayload {
			stepRequestResponsePayload[strings.ReplaceAll(s, " ", "_")] = map[string]interface{}{"request": stepRequestResponse.Request, "response": stepRequestResponse.Response}
		}
		step.canStepExecute, err = hooks.GetExprHook().ShouldStepExecute(step.When, stepRequestResponsePayload, prefix)
	}

	return err
}

func (step *Step) stepExecution(requestBody *StepRequest, prefix string) (interface{}, error) {
	switch step.Mode {
	case "HTTP":
		step.UpdateRequestHeadersBasedOnRequestHeadersAndStepHeaders(requestBody)
		res, err := step.Val.(*executors.HTTPVal).DoExecute(requestBody.Payload, prefix)
		return res, err
	case "AMQP":
		res, err := step.Val.(*executors.AMQPVal).DoExecute(requestBody, prefix)
		return res, err
	case "KAFKA":
		res, err := step.Val.(*executors.KafkaVal).DoExecute(requestBody, prefix)
		return res, err
	}
	panic("Invalid mode specified")
}

// UpdateRequestHeadersBasedOnRequestHeadersAndStepHeaders updates request headers based on request and step headers
func (step *Step) UpdateRequestHeadersBasedOnRequestHeadersAndStepHeaders(requestBody *StepRequest) {
	headers := step.Val.(*executors.HTTPVal).Headers
	requestHeaders := requestBody.Headers
	if requestHeaders != "" && headers != "" {
		headers = requestHeaders + headers
	} else if requestHeaders != "" && headers == "" {
		headers = requestHeaders
	}
	step.Val.(*executors.HTTPVal).Headers = headers
}

// DoExecute executes the step
func (step *Step) DoExecute(requestContext RequestContext, prefix string) (_ interface{}, _ error) {
	err := step.PreStepExecution(requestContext.StepsContext, prefix)
	if err != nil {
		return nil, err
	}
	request := requestContext.GetStepRequestFromContext(step.Name)
	if !step.canStepExecute {
		requestContext.StepsContext[step.Name].StepSkipped = true
		log.Printf("%s Skipping step: %s, condition (%s), request payload (%v), not satisified ", prefix, step.Name, step.When, requestContext.StepsContext)
		return request, nil
	}
	res, err := step.stepExecution(NewStepRequest(requestContext.ServiceRequestID, step.ID, request, requestContext.GetStepRequestHeadersFromContext(step.Name)), prefix)
	//post Step execution
	return res, err
}

// DoTransform transforms the request
func (step *Step) DoTransform(requestContext RequestContext, prefix string) (map[string]interface{}, error) {
	stepRequestResponsePayload := make(map[string]interface{})
	for s, stepRequestResponse := range requestContext.StepsContext {
		stepRequestResponsePayload[strings.ReplaceAll(s, " ", "_")] = map[string]interface{}{"request": stepRequestResponse.Request, "response": stepRequestResponse.Response}
	}
	if step.Transform {
		switch step.TransformFormat {
		case "XML":
			res, err := step.RequestTransform.(*transform.XMLTransform).DoTransform(stepRequestResponsePayload, prefix)
			return res, err
		default:
			res, err := step.RequestTransform.(*transform.JSONTransform).DoTransform(stepRequestResponsePayload, prefix)
			return res, err
		}
	}
	return stepRequestResponsePayload, nil
}

// UnmarshalJSON unmarshals step from JSON
func (step *Step) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	mode := v["mode"]
	//TODO I guess this is initialization section otherwise transform.JSONTransform was not getting called.
	requestTransform := v["transformFormat"]

	if requestTransform != nil {
		transformErr := step.setRequestTransform(requestTransform)
		if transformErr != nil {
			return transformErr
		}
	}
	err := step.setMode(mode)
	if err != nil {
		return err
	}
	type stepStruct Step
	err = json.Unmarshal(data, (*stepStruct)(step))
	return err
}

func (step *Step) setRequestTransform(requestTransform interface{}) error {
	m, ok := requestTransform.(string)
	if !ok {
		return fmt.Errorf("%s is an invalid Mode", requestTransform)
	}
	switch m {
	case "XML":
		step.RequestTransform = &transform.XMLTransform{}
	default:
		step.RequestTransform = &transform.JSONTransform{}
	}
	return nil
}

func (step *Step) setMode(mode interface{}) error {
	m, ok := mode.(string)
	if !ok {
		return fmt.Errorf("%s is an invalid Mode", mode)
	}
	switch m {
	case "HTTP":
		step.Val = &executors.HTTPVal{}
	case "AMQP":
		step.Val = &executors.AMQPVal{}
	case "KAFKA":
		step.Val = &executors.KafkaVal{}
	default:
		return fmt.Errorf("%s is an invalid Mode", mode)
	}
	return nil
}

func (step Step) getHTTPVal() executors.HTTPVal {
	return step.Val.(executors.HTTPVal)
}

func (step Step) getAMQPVal() *executors.AMQPVal {
	return step.Val.(*executors.AMQPVal)
}
