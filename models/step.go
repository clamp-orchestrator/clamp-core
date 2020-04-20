package models

import (
	"clamp-core/executors"
	"encoding/json"
	"fmt"
)

type Val interface {
}

type Step struct {
	Id        int    `json:"id"`
	Name      string `json:"name" binding:"required"`
	StepType  string `json:"type" binding:"required,oneof=SYNC ASYNC"`
	Mode      string `json:"mode" binding:"required,oneof=HTTP AMQP"`
	Val       Val    `json:"val" binding:"required"`
	Transform bool   `json:"transform"`
	Enabled   bool   `json:"enabled"`
}

func (step Step) DoExecute(requestBody interface{}) (interface{}, error) {
	switch step.Mode {
	case "HTTP":
		return step.Val.(*executors.HttpVal).DoExecute(requestBody)
	case "AMQP":
		return step.Val.(*executors.AMQPVal).DoExecute(requestBody)
	}
	panic("Invalid mode specified")
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
