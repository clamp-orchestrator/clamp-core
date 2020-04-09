package models

import (
	"clamp-core/executors"
	"encoding/json"
	"fmt"
	"log"
)

type Val interface {
}

type Step struct {
	Id        string `json:"id"`
	Name      string `json:"name" binding:"required"`
	Mode      string `json:"mode" binding:"required,oneof=HTTP QUEUE"`
	Val       Val    `json:"val" binding:"required"`
	Transform bool   `json:"transform"`
	Enabled   bool   `json:"enabled"`
}

func (step Step) DoExecute(requestBody interface{}) (interface{}, error) {
	switch step.Mode {
	case "HTTP":
		log.Println("Inside HTTP Execute")
		return step.Val.(*executors.HttpVal).DoExecute(requestBody)
	case "QUEUE":
		log.Println("Inside QUEUE Execute")
		return step.Val.(*executors.QueueVal).DoExecute()
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
	case "QUEUE":
		step.Val = &executors.QueueVal{}
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
