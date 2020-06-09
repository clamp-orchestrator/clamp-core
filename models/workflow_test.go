package models

import (
	"clamp-core/config"
	"clamp-core/executors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldCreateANewWorkflow(t *testing.T) {
	http := executors.HttpVal{
		Method:  "GET",
		Url:     "http://18.236.212.57:3333/api/v1/user",
		Headers: "",
	}
	steps := []Step{{}}
	steps[0] = Step{
		Name:    "firstStep",
		Type:    "SYNC",
		Mode:    "HTTP",
		Val:     http,
		Enabled: true,
	}
	workflow := Workflow{
		Id:          "1",
		Name:        "Test",
		Description: "Test",
		Enabled:     false,
		Steps:       steps,
	}

	serviceFlowRequest := workflow
	err := binding.Validator.ValidateStruct(workflow)
	if err != nil {
		log.Println(err)
	}
	assert.Nil(t, err)
	workflowResponse := CreateWorkflow(serviceFlowRequest)

	assert.NotEmpty(t, workflowResponse.Id)
	assert.NotNil(t, workflowResponse.CreatedAt)
	assert.Equal(t, serviceFlowRequest.Description, workflowResponse.Description, fmt.Sprintf("Expected workflow description to be %s but was %s", serviceFlowRequest.Description, workflowResponse.Description))
	assert.Equal(t, serviceFlowRequest.Name, workflowResponse.Name, fmt.Sprintf("Expected worflow name to be %s but was %s", serviceFlowRequest.Name, workflowResponse.Name))
	assert.Equal(t, serviceFlowRequest.Steps[0].Name, workflowResponse.Steps[0].Name, fmt.Sprintf("Expected worflow first step name to be %s but was %s", serviceFlowRequest.Steps[0].Name, workflowResponse.Steps[0].Name))
	assert.Equal(t, "GET", workflowResponse.Steps[0].getHttpVal().Method)
	assert.Equal(t, "http://18.236.212.57:3333/api/v1/user", workflowResponse.Steps[0].getHttpVal().Url)
	assert.Equal(t, "", workflowResponse.Steps[0].getHttpVal().Headers)
}

func TestShouldNotCreateWorkflowIfStepValIsNotPresent(t *testing.T) {
	steps := []Step{{}}
	steps[0] = Step{
		Name:    "firstStep",
		Type:    "SYNC",
		Mode:    "HTTP",
		Enabled: true,
	}
	workflow := Workflow{
		Id:          "1",
		Name:        "Test",
		Description: "Test",
		Enabled:     false,
		Steps:       steps,
	}

	err := binding.Validator.ValidateStruct(workflow)
	if err != nil {
		log.Println(err)
	}
	assert.NotNil(t, err)
	assert.Equal(t, "Key: 'Workflow.Steps[0].Val' Error:Field validation for 'Val' failed on the 'required' tag", err.Error())
}

func TestShouldThrowErrorIfInvalidModeIsUsed(t *testing.T) {
	http := executors.HttpVal{
		Method:  "GET",
		Url:     "http://18.236.212.57:3333/api/v1/user",
		Headers: "",
	}
	steps := []Step{{}}
	const InvalidMode = "xyz"
	steps[0] = Step{
		Name:    "firstStep",
		Type:    "SYNC",
		Mode:    InvalidMode,
		Val:     http,
		Enabled: true,
	}
	workflow := Workflow{
		Id:          "1",
		Name:        "Test",
		Description: "Test",
		Enabled:     false,
		Steps:       steps,
	}

	err := binding.Validator.ValidateStruct(workflow)
	if err != nil {
		log.Println(err)
	}
	assert.NotNil(t, err)
	assert.Equal(t, "Key: 'Workflow.Steps[0].Mode' Error:Field validation for 'Mode' failed on the 'oneof' tag", err.Error())
}

func TestShouldThrowErrorIfGetHTTPValIsCalledForADiffMode(t *testing.T) {
	queue := executors.AMQPVal{
		ConnectionURL: "http://18.236.212.57:3333",
		QueueName:     "topic-a",
	}
	steps := []Step{{}}
	steps[0] = Step{
		Name:    "firstStep",
		Type:    "ASYNC",
		Mode:    "AMQP",
		Val:     queue,
		Enabled: true,
	}
	workflow := Workflow{
		Id:          "1",
		Name:        "Test",
		Description: "Test",
		Enabled:     false,
		Steps:       steps,
	}

	err := binding.Validator.ValidateStruct(workflow)
	assert.Nil(t, err)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetHttpVal should have panicked!")
		}
	}()
	workflow.Steps[0].getHttpVal()
}

func TestShouldThrowErrorIfGetHTTPValUrlIsEmpty(t *testing.T) {
	http := executors.HttpVal{
		Method:  "GET",
		Url:     "",
		Headers: "",
	}
	steps := []Step{{}}
	steps[0] = Step{
		Name:    "firstStep",
		Type:    "ASYNC",
		Mode:    "AMQP",
		Val:     http,
		Enabled: true,
	}
	workflow := Workflow{
		Id:          "1",
		Name:        "Test",
		Description: "Test",
		Enabled:     false,
		Steps:       steps,
	}

	err := binding.Validator.ValidateStruct(workflow)
	assert.NotNil(t, err)
	assert.Equal(t, "Key: 'Workflow.Steps[0].Val.Url' Error:Field validation for 'Url' failed on the 'required' tag", err.Error())
}

func TestIfReplyToQueueNameIsNotProvidedAsPartOfWorkflowRequestShouldReadDefaultValueFromConfig(t *testing.T) {
	queue := &executors.AMQPVal{
		ConnectionURL: "http://18.236.212.57:3333",
		QueueName:     "topic-a",
	}
	steps := []Step{{}}
	steps[0] = Step{
		Name:    "firstStep",
		Type:    "ASYNC",
		Mode:    "AMQP",
		Val:     queue,
		Enabled: true,
	}
	workflow := Workflow{
		Id:          "1",
		Name:        "Test",
		Description: "Test",
		Enabled:     false,
		Steps:       steps,
	}

	serviceFlowRequest := workflow
	err := binding.Validator.ValidateStruct(workflow)
	if err != nil {
		log.Println(err)
	}
	assert.Nil(t, err)
	workflowResponse := CreateWorkflow(serviceFlowRequest)

	assert.NotEmpty(t, workflowResponse.Id)
	assert.NotNil(t, workflowResponse.CreatedAt)
	assert.Equal(t, serviceFlowRequest.Description, workflowResponse.Description, fmt.Sprintf("Expected workflow description to be %s but was %s", serviceFlowRequest.Description, workflowResponse.Description))
	assert.Equal(t, serviceFlowRequest.Name, workflowResponse.Name, fmt.Sprintf("Expected worflow name to be %s but was %s", serviceFlowRequest.Name, workflowResponse.Name))
	assert.Equal(t, serviceFlowRequest.Steps[0].Name, workflowResponse.Steps[0].Name, fmt.Sprintf("Expected worflow first step name to be %s but was %s", serviceFlowRequest.Steps[0].Name, workflowResponse.Steps[0].Name))
	assert.Equal(t, "http://18.236.212.57:3333", workflowResponse.Steps[0].getAMQPVal().ConnectionURL)
	assert.Equal(t, config.ENV.QueueName, workflowResponse.Steps[0].getAMQPVal().ReplyTo)
	assert.Equal(t, "topic-a", workflowResponse.Steps[0].getAMQPVal().QueueName)
}
