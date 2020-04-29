package services

import (
	"clamp-core/executors"
	"clamp-core/models"
	"clamp-core/transform"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const workflowName string = "testWF"

func TestAddServiceRequestToChannel(t *testing.T) {
	var functionCalledStack []string
	findWorkflowByNameMock = func(workflowName string) (workflow models.Workflow, err error) {
		workflow.Id = "TEST_WF"
		step := models.Step{
			Name:      "1",
			StepType:  "SYNC",
			Mode:      "HTTP",
			Transform: false,
			Enabled:   false,
			Val: &executors.HttpVal{
				Method:  "POST",
				Url:     "http://34.222.238.234:3333/api/v1/login",
				Headers: "",
			},
		}
		workflow.Steps = []models.Step{step}
		functionCalledStack = append(functionCalledStack, "findWorkflowByName")
		return workflow, err
	}
	saveStepStatusMock = func(stepStatus models.StepsStatus) (status models.StepsStatus, err error) {
		functionCalledStack = append(functionCalledStack, "saveStepStatusMock")
		return status, nil
	}
	type args struct {
		serviceReq models.ServiceRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should process service request",
			args: args{
				serviceReq: models.ServiceRequest{
					ID:           uuid.New(),
					WorkflowName: workflowName,
					Status:       models.STATUS_NEW,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddServiceRequestToChannel(tt.args.serviceReq)
			time.Sleep(time.Second * 5)
			assert.Equal(t, 0, len(serviceRequestChannel))
			assert.Equal(t, 3, len(functionCalledStack))
		})
	}
}

func AddServiceRequestToChannelWithTransformationEnabledForOneStepInTheWorkflowTest(t *testing.T) {
	var functionCalledStack []string
	findWorkflowByNameMock = func(workflowName string) (workflow models.Workflow, err error) {
		workflow.Id = "TEST_WF"
		step := models.Step{
			Name:      "1",
			StepType:  "SYNC",
			Mode:      "HTTP",
			Transform: true,
			Enabled:   false,
			RequestTransform: &transform.JsonTransform{
				Spec:map[string]interface{}{"name":"test"},
			},
			Val: &executors.HttpVal{
				Method:  "POST",
				Url:     "http://34.222.238.234:3333/api/v1/login",
				Headers: "",
			},
		}
		workflow.Steps = []models.Step{step}
		functionCalledStack = append(functionCalledStack, "findWorkflowByName")
		return workflow, err
	}
	saveStepStatusMock = func(stepStatus models.StepsStatus) (status models.StepsStatus, err error) {
		functionCalledStack = append(functionCalledStack, "saveStepStatusMock")
		return status, nil
	}
	type args struct {
		serviceReq models.ServiceRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should process service request",
			args: args{
				serviceReq: models.ServiceRequest{
					ID:           uuid.New(),
					WorkflowName: workflowName,
					Status:       models.STATUS_NEW,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddServiceRequestToChannel(tt.args.serviceReq)
			time.Sleep(time.Second * 5)
			assert.Equal(t, 0, len(serviceRequestChannel))
			assert.Equal(t, 3, len(functionCalledStack))
		})
	}
}



