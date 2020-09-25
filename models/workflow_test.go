package models

import (
	"clamp-core/config"
	"clamp-core/executors"
	"fmt"
	"log"
	"testing"

	"github.com/gin-gonic/gin/binding"

	"github.com/stretchr/testify/assert"
)

func TestShouldCreateANewWorkflow(t *testing.T) {
	http := executors.HTTPVal{
		Method:  "GET",
		URL:     "https://run.mocky.io/v3/0590fbf8-0f1c-401c-b9df-65e98ef0385d",
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
		ID:          "1",
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

	assert.NotEmpty(t, workflowResponse.ID)
	assert.NotNil(t, workflowResponse.CreatedAt)
	assert.Equal(t, serviceFlowRequest.Description, workflowResponse.Description, fmt.Sprintf("Expected workflow description to be %s but was %s", serviceFlowRequest.Description, workflowResponse.Description))
	assert.Equal(t, serviceFlowRequest.Name, workflowResponse.Name, fmt.Sprintf("Expected worflow name to be %s but was %s", serviceFlowRequest.Name, workflowResponse.Name))
	assert.Equal(t, serviceFlowRequest.Steps[0].Name, workflowResponse.Steps[0].Name, fmt.Sprintf("Expected worflow first step name to be %s but was %s", serviceFlowRequest.Steps[0].Name, workflowResponse.Steps[0].Name))
	assert.Equal(t, "GET", workflowResponse.Steps[0].getHTTPVal().Method)
	assert.Equal(t, "https://run.mocky.io/v3/0590fbf8-0f1c-401c-b9df-65e98ef0385d", workflowResponse.Steps[0].getHTTPVal().URL)
	assert.Equal(t, "", workflowResponse.Steps[0].getHTTPVal().Headers)
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
		ID:          "1",
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
	http := executors.HTTPVal{
		Method:  "GET",
		URL:     "https://run.mocky.io/v3/0590fbf8-0f1c-401c-b9df-65e98ef0385d",
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
		ID:          "1",
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
		ConnectionURL: "https://run.mocky.io/v3/0590fbf8-0f1c-401c-b9df-65e98ef0385d",
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
		ID:          "1",
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
	workflow.Steps[0].getHTTPVal()
}

func TestShouldThrowErrorIfGetHTTPValUrlIsEmpty(t *testing.T) {
	http := executors.HTTPVal{
		Method:  "GET",
		URL:     "",
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
		ID:          "1",
		Name:        "Test",
		Description: "Test",
		Enabled:     false,
		Steps:       steps,
	}

	err := binding.Validator.ValidateStruct(workflow)
	assert.NotNil(t, err)
	assert.Equal(t, "Key: 'Workflow.Steps[0].Val.URL' Error:Field validation for 'URL' failed on the 'required' tag", err.Error())
}

func TestIfReplyToQueueNameIsNotProvidedAsPartOfWorkflowRequestShouldReadDefaultValueFromConfig(t *testing.T) {
	queue := &executors.AMQPVal{
		ConnectionURL: "https://run.mocky.io/v3/0590fbf8-0f1c-401c-b9df-65e98ef0385d",
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
		ID:          "1",
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

	assert.NotEmpty(t, workflowResponse.ID)
	assert.NotNil(t, workflowResponse.CreatedAt)
	assert.Equal(t, serviceFlowRequest.Description, workflowResponse.Description, fmt.Sprintf("Expected workflow description to be %s but was %s", serviceFlowRequest.Description, workflowResponse.Description))
	assert.Equal(t, serviceFlowRequest.Name, workflowResponse.Name, fmt.Sprintf("Expected worflow name to be %s but was %s", serviceFlowRequest.Name, workflowResponse.Name))
	assert.Equal(t, serviceFlowRequest.Steps[0].Name, workflowResponse.Steps[0].Name, fmt.Sprintf("Expected worflow first step name to be %s but was %s", serviceFlowRequest.Steps[0].Name, workflowResponse.Steps[0].Name))
	assert.Equal(t, "https://run.mocky.io/v3/0590fbf8-0f1c-401c-b9df-65e98ef0385d", workflowResponse.Steps[0].getAMQPVal().ConnectionURL)
	assert.Equal(t, config.ENV.QueueName, workflowResponse.Steps[0].getAMQPVal().ReplyTo)
	assert.Equal(t, "topic-a", workflowResponse.Steps[0].getAMQPVal().QueueName)
}

func TestShouldCreateNewWorkflowWithOnFailureSteps(t *testing.T) {
	http := executors.HTTPVal{
		Method:  "GET",
		URL:     "https://run.mocky.io/v3/0590fbf8-0f1c-401c-b9df-65e98ef0385d",
		Headers: "",
	}
	steps := []Step{{}}
	failureSteps := []Step{{}}
	failureSteps[0] = Step{
		Name:    "onFailureStep",
		Mode:    "HTTP",
		Val:     http,
		Enabled: true,
	}
	steps[0] = Step{
		Name:      "firstStep",
		Mode:      "HTTP",
		Val:       http,
		Enabled:   true,
		OnFailure: failureSteps,
	}
	workflow := Workflow{
		ID:          "1",
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

	assert.NotEmpty(t, workflowResponse.ID)
	assert.NotNil(t, workflowResponse.CreatedAt)
	assert.Equal(t, serviceFlowRequest.Description, workflowResponse.Description, fmt.Sprintf("Expected workflow description to be %s but was %s", serviceFlowRequest.Description, workflowResponse.Description))
	assert.Equal(t, serviceFlowRequest.Name, workflowResponse.Name, fmt.Sprintf("Expected worflow name to be %s but was %s", serviceFlowRequest.Name, workflowResponse.Name))
	assert.Equal(t, serviceFlowRequest.Steps[0].Name, workflowResponse.Steps[0].Name, fmt.Sprintf("Expected worflow first step name to be %s but was %s", serviceFlowRequest.Steps[0].Name, workflowResponse.Steps[0].Name))
	assert.Equal(t, "https://run.mocky.io/v3/0590fbf8-0f1c-401c-b9df-65e98ef0385d", workflowResponse.Steps[0].getHTTPVal().URL)
	assert.NotNil(t, workflowResponse.Steps[0].OnFailure)
	assert.Equal(t, "onFailureStep", workflowResponse.Steps[0].OnFailure[0].Name)
	assert.Equal(t, "HTTP", workflowResponse.Steps[0].OnFailure[0].Mode)
	assert.Equal(t, "https://run.mocky.io/v3/0590fbf8-0f1c-401c-b9df-65e98ef0385d", workflowResponse.Steps[0].OnFailure[0].getHTTPVal().URL)
}
