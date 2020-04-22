package models

import (
	"clamp-core/executors"
	"clamp-core/hooks"
	"encoding/json"
	"fmt"
	"log"
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
}

func (step *Step) CanStepExecute(canStepExecute bool) {
	step.canStepExecute = canStepExecute
}

func (step *Step) preStepExecution(requestBody StepRequest, prefix string) error {
	canStepExecute, err := hooks.PreStepHookExecutor(step.When, requestBody.Payload, prefix)
	step.canStepExecute = canStepExecute
	return err
}

func (step *Step) DoExecute(requestBody StepRequest, prefix string) (skipStepExecution bool, _ interface{}, _ error) {
	err := step.preStepExecution(requestBody, prefix)
	if err != nil {
		return skipStepExecution, nil, err
	}
	if skipStepExecution = !step.canStepExecute; !skipStepExecution {
		switch step.Mode {
		case "HTTP":
			res, err := step.Val.(*executors.HttpVal).DoExecute(requestBody, prefix)
			return skipStepExecution, res, err
		case "AMQP":
			res, err := step.Val.(*executors.AMQPVal).DoExecute(requestBody, prefix)
			return skipStepExecution, res, err
		}
		panic("Invalid mode specified")
	} else {
		log.Printf("%s Skipping step: %s, condition (%s), request payload (%v), not satisified ", prefix, step.Name, step.When, requestBody.Payload)
		return skipStepExecution, requestBody, nil
	}
}

func (step *Step) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	mode := v["mode"]
	switch mode {
	case "HTTP":
		step.Val = &executors.HttpVal{}
	case "AMQP":
		step.Val = &executors.AMQPVal{}
	default:
		return fmt.Errorf("%s is an invalid Mode", mode)
	}
	type stepStruct Step
	err := json.Unmarshal(data, (*stepStruct)(step))
	return err
}

func (step Step) getHttpVal() executors.HttpVal {
	return step.Val.(executors.HttpVal)
}
