package models

import (
	"clamp-core/config"
	"clamp-core/executors"
	"clamp-core/utils"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/gin-gonic/gin/binding"

	"github.com/stretchr/testify/assert"
)

func TestShouldCreateANewWorkflow(t *testing.T) {
	http := executors.HTTPVal{
		Method:  "GET",
		URL:     testHTTPServer.URL,
		Headers: "",
	}
	steps := []Step{{}}
	steps[0] = Step{
		Name:    "firstStep",
		Type:    utils.StepTypeSync,
		Mode:    utils.StepModeHTTP,
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
	workflowResponse := CreateWorkflow(&serviceFlowRequest)

	assert.NotEmpty(t, workflowResponse.ID)
	assert.NotNil(t, workflowResponse.CreatedAt)
	assert.Equal(t, serviceFlowRequest.Description, workflowResponse.Description, fmt.Sprintf("Expected workflow description to be %s but was %s", serviceFlowRequest.Description, workflowResponse.Description))
	assert.Equal(t, serviceFlowRequest.Name, workflowResponse.Name, fmt.Sprintf("Expected worflow name to be %s but was %s", serviceFlowRequest.Name, workflowResponse.Name))
	assert.Equal(t, serviceFlowRequest.Steps[0].Name, workflowResponse.Steps[0].Name, fmt.Sprintf("Expected worflow first step name to be %s but was %s", serviceFlowRequest.Steps[0].Name, workflowResponse.Steps[0].Name))
	assert.Equal(t, "GET", workflowResponse.Steps[0].getHTTPVal().Method)
	assert.Equal(t, testHTTPServer.URL, workflowResponse.Steps[0].getHTTPVal().URL)
	assert.Equal(t, "", workflowResponse.Steps[0].getHTTPVal().Headers)
}

func TestShouldNotCreateWorkflowIfStepValIsNotPresent(t *testing.T) {
	steps := []Step{{}}
	steps[0] = Step{
		Name:    "firstStep",
		Type:    utils.StepTypeSync,
		Mode:    utils.StepModeHTTP,
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
		URL:     testHTTPServer.URL,
		Headers: "",
	}
	steps := []Step{{}}
	const InvalidMode = "xyz"
	steps[0] = Step{
		Name:    "firstStep",
		Type:    utils.StepTypeSync,
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
		ConnectionURL: testHTTPServer.URL,
		QueueName:     "topic-a",
	}
	steps := []Step{{}}
	steps[0] = Step{
		Name:    "firstStep",
		Type:    utils.StepTypeAsync,
		Mode:    utils.StepModeAMQP,
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
		Type:    utils.StepTypeAsync,
		Mode:    utils.StepModeAMQP,
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
		ConnectionURL: testHTTPServer.URL,
		QueueName:     "topic-a",
	}
	steps := []Step{{}}
	steps[0] = Step{
		Name:    "firstStep",
		Type:    utils.StepTypeAsync,
		Mode:    utils.StepModeAMQP,
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
	workflowResponse := CreateWorkflow(&serviceFlowRequest)

	assert.NotEmpty(t, workflowResponse.ID)
	assert.NotNil(t, workflowResponse.CreatedAt)
	assert.Equal(t, serviceFlowRequest.Description, workflowResponse.Description, fmt.Sprintf("Expected workflow description to be %s but was %s", serviceFlowRequest.Description, workflowResponse.Description))
	assert.Equal(t, serviceFlowRequest.Name, workflowResponse.Name, fmt.Sprintf("Expected worflow name to be %s but was %s", serviceFlowRequest.Name, workflowResponse.Name))
	assert.Equal(t, serviceFlowRequest.Steps[0].Name, workflowResponse.Steps[0].Name, fmt.Sprintf("Expected worflow first step name to be %s but was %s", serviceFlowRequest.Steps[0].Name, workflowResponse.Steps[0].Name))
	assert.Equal(t, testHTTPServer.URL, workflowResponse.Steps[0].getAMQPVal().ConnectionURL)
	assert.Equal(t, config.ENV.QueueName, workflowResponse.Steps[0].getAMQPVal().ReplyTo)
	assert.Equal(t, "topic-a", workflowResponse.Steps[0].getAMQPVal().QueueName)
}

func TestShouldCreateNewWorkflowWithOnFailureSteps(t *testing.T) {
	http := executors.HTTPVal{
		Method:  "GET",
		URL:     testHTTPServer.URL,
		Headers: "",
	}
	steps := []Step{{}}
	failureSteps := []Step{{}}
	failureSteps[0] = Step{
		Name:    "onFailureStep",
		Mode:    utils.StepModeHTTP,
		Val:     http,
		Enabled: true,
	}
	steps[0] = Step{
		Name:      "firstStep",
		Mode:      utils.StepModeHTTP,
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
	workflowResponse := CreateWorkflow(&serviceFlowRequest)

	assert.NotEmpty(t, workflowResponse.ID)
	assert.NotNil(t, workflowResponse.CreatedAt)
	assert.Equal(t, serviceFlowRequest.Description, workflowResponse.Description, fmt.Sprintf("Expected workflow description to be %s but was %s", serviceFlowRequest.Description, workflowResponse.Description))
	assert.Equal(t, serviceFlowRequest.Name, workflowResponse.Name, fmt.Sprintf("Expected worflow name to be %s but was %s", serviceFlowRequest.Name, workflowResponse.Name))
	assert.Equal(t, serviceFlowRequest.Steps[0].Name, workflowResponse.Steps[0].Name, fmt.Sprintf("Expected worflow first step name to be %s but was %s", serviceFlowRequest.Steps[0].Name, workflowResponse.Steps[0].Name))
	assert.Equal(t, testHTTPServer.URL, workflowResponse.Steps[0].getHTTPVal().URL)
	assert.NotNil(t, workflowResponse.Steps[0].OnFailure)
	assert.Equal(t, "onFailureStep", workflowResponse.Steps[0].OnFailure[0].Name)
	assert.Equal(t, utils.StepModeHTTP, workflowResponse.Steps[0].OnFailure[0].Mode)
	assert.Equal(t, testHTTPServer.URL, workflowResponse.Steps[0].OnFailure[0].getHTTPVal().URL)
}

func TestWorkflowStepIDAreSequentiallyIncrementing(t *testing.T) {
	w := CreateWorkflow(&Workflow{
		ID:          "abc",
		Name:        "WorkflowABC",
		Description: "Workflow ABC",
		Enabled:     true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Steps: []Step{
			{
				ID:             0,
				Name:           "Step1",
				Mode:           utils.StepModeHTTP,
				Enabled:        true,
				canStepExecute: true,
				OnFailure: []Step{
					{
						ID:             0,
						Name:           "OnFailureStep1",
						Mode:           utils.StepModeHTTP,
						Enabled:        true,
						canStepExecute: true,
					},
					{
						ID:             0,
						Name:           "OnFailureStep2",
						Mode:           utils.StepModeHTTP,
						Enabled:        true,
						canStepExecute: true,
					},
				},
			},
			{
				ID:             0,
				Name:           "Step2",
				Mode:           utils.StepModeHTTP,
				Enabled:        true,
				canStepExecute: true,
				OnFailure: []Step{
					{
						ID:             0,
						Name:           "OnFailureStep1",
						Mode:           utils.StepModeHTTP,
						Enabled:        true,
						canStepExecute: true,
					},
					{
						ID:             0,
						Name:           "OnFailureStep2",
						Mode:           utils.StepModeHTTP,
						Enabled:        true,
						canStepExecute: true,
					},
				},
			},
			{
				ID:             0,
				Name:           "Step2",
				Mode:           utils.StepModeHTTP,
				Enabled:        true,
				canStepExecute: true,
			},
		},
	})

	assert.Equal(t, 1, w.Steps[0].ID)
	assert.Equal(t, 2, w.Steps[0].OnFailure[0].ID)
	assert.Equal(t, 3, w.Steps[0].OnFailure[1].ID)

	assert.Equal(t, 4, w.Steps[1].ID)
	assert.Equal(t, 5, w.Steps[1].OnFailure[0].ID)
	assert.Equal(t, 6, w.Steps[1].OnFailure[1].ID)

	assert.Equal(t, 7, w.Steps[2].ID)
}
