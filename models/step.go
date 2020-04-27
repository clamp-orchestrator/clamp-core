package models

import (
	"clamp-core/executors"
	"clamp-core/hooks"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Val interface {
}

type Step struct {
	Id             int    `json:"id"`
	Name           string `json:"name" binding:"required"`
	StepType       string `json:"type" binding:"required,oneof=SYNC ASYNC"`
	Mode           string `json:"mode" binding:"required,oneof=HTTP AMQP"`
	Val            Val    `json:"val" binding:"required"`
	Transform      bool   `json:"transform"`
	Enabled        bool   `json:"enabled"`
	When           string `json:"when"`
	canStepExecute bool
	//shouldStepExecute func(whenCondition string, stepRequest map[string]interface{}, prefix string) (canStepExecute bool, _ error)
	//transformRequest  func(stepRequest map[string]interface{}, prefix string) (map[string]interface{}, error)
}

func (step *Step) DidStepExecute() bool {
	return step.canStepExecute
}

func (step *Step) preStepExecution(contextPayload map[string]RequestResponse, prefix string) (err error) {
	step.canStepExecute = true
	stepRequestResponsePayload := make(map[string]interface{})

	if step.When != "" {
		for s, stepRequestResponse := range contextPayload {
			stepRequestResponsePayload[strings.ReplaceAll(s," ", "_")] = map[string]interface{}{"request": stepRequestResponse.Request, "response":stepRequestResponse.Response}
		}
		step.canStepExecute, err = hooks.GetExprHook().ShouldStepExecute(step.When, stepRequestResponsePayload, prefix)
	}

	return err
}

func (step *Step) stepExecution(requestBody StepRequest, prefix string) (interface{}, error) {
	switch step.Mode {
	case "HTTP":
		res, err := step.Val.(*executors.HttpVal).DoExecute(requestBody, prefix)
		return res, err
	case "AMQP":
		res, err := step.Val.(*executors.AMQPVal).DoExecute(requestBody, prefix)
		return res, err
	}
	panic("Invalid mode specified")
}

func (step *Step) DoExecute(requestBody StepRequest, prefix string, requestContext RequestContext) (_ interface{}, _ error) {
	err := step.preStepExecution(requestContext.Payload, prefix)
	if err != nil {
		return nil, err
	}
	if !step.canStepExecute {
		log.Printf("%s Skipping step: %s, condition (%s), request payload (%v), not satisified ", prefix, step.Name, step.When, requestContext.Payload)
		return requestBody, nil
	}
	res, err := step.stepExecution(requestBody, prefix)
	//post Step execution
	return res, err
}

func (step *Step) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	mode := v["mode"]
	err := step.setMode(mode)
	if err != nil {
		return err
	}
	type stepStruct Step
	err = json.Unmarshal(data, (*stepStruct)(step))
	return err
}

func (step *Step) setMode(mode interface{}) error {
	m, ok := mode.(string)
	if !ok {
		return fmt.Errorf("%s is an invalid Mode", mode)
	}
	switch m {
	case "HTTP":
		step.Val = &executors.HttpVal{}
	case "AMQP":
		step.Val = &executors.AMQPVal{}
	default:
		return fmt.Errorf("%s is an invalid Mode", mode)
	}
	return nil
}

func (step Step) getHttpVal() executors.HttpVal {
	return step.Val.(executors.HttpVal)
}
