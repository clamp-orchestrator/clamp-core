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
		workflow.ID = "TEST_WF"
		step := models.Step{
			Name:      "1",
			Type:      "SYNC",
			Mode:      "HTTP",
			Transform: false,
			Enabled:   false,
			Val: &executors.HTTPVal{
				Method:  "POST",
				URL:     "http://54.190.25.178:3333/api/v1/login",
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

func TestShouldAddServiceRequestToChannelWithTransformationEnabledForOneStepInTheWorkflow(t *testing.T) {
	var functionCalledStack []string
	findWorkflowByNameMock = func(workflowName string) (workflow models.Workflow, err error) {
		workflow.ID = "TEST_WF"
		step := models.Step{
			Name:      "1",
			Type:      "SYNC",
			Mode:      "HTTP",
			Transform: true,
			Enabled:   false,
			RequestTransform: &transform.JSONTransform{
				Spec: map[string]interface{}{"name": "test"},
			},
			Val: &executors.HTTPVal{
				Method:  "POST",
				URL:     "http://54.190.25.178:3333/api/v1/login",
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

func TestShouldSkipStepIfConditionDoesNotMatch(t *testing.T) {
	var functionCalledStack []string
	findWorkflowByNameMock = func(workflowName string) (workflow models.Workflow, err error) {
		workflow.ID = "TEST_WF"
		step := models.Step{
			Name:      "skipStep",
			Type:      "SYNC",
			Mode:      "HTTP",
			Transform: false,
			Enabled:   false,
			When:      "skipStep.request.id1 == 'val3'",
			Val: &executors.HTTPVal{
				Method:  "POST",
				URL:     "http://54.190.25.178:3333/api/v1/login",
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
			name: "Should skip the step when condition does not match with service request payload",
			args: args{
				serviceReq: models.ServiceRequest{
					ID:           uuid.New(),
					WorkflowName: workflowName,
					Status:       models.STATUS_NEW,
					Payload:      prepareRequestPayload(),
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

func TestShouldResumeTheWorkflowExecutionFromNextStep(t *testing.T) {
	var functionCalledStack []string
	findWorkflowByNameMock = func(workflowName string) (workflow models.Workflow, err error) {
		workflow.ID = "TEST_WF"
		step := models.Step{
			Name:      "firstStep",
			Type:      "SYNC",
			Mode:      "HTTP",
			Transform: false,
			Enabled:   false,
			When:      "firstStep.request.id1 == 'val1'",
			Val: &executors.HTTPVal{
				Method:  "POST",
				URL:     "http://54.190.25.178:3333/api/v1/login",
				Headers: "",
			},
		}
		step1 := models.Step{
			Name:      "secondStep",
			Type:      "SYNC",
			Mode:      "HTTP",
			Transform: false,
			Enabled:   false,
			When:      "firstStep.request.id1 == 'val1'",
			Val: &executors.HTTPVal{
				Method:  "POST",
				URL:     "http://54.190.25.178:3333/api/v1/login",
				Headers: "",
			},
		}
		workflow.Steps = []models.Step{step, step1}
		functionCalledStack = append(functionCalledStack, "findWorkflowByName")
		return workflow, err
	}
	saveStepStatusMock = func(stepStatus models.StepsStatus) (status models.StepsStatus, err error) {
		functionCalledStack = append(functionCalledStack, "saveStepStatusMock")
		return status, nil
	}
	findStepStatusByServiceRequestIDAndStatusMock = func(serviceRequestId uuid.UUID, status models.Status) ([]models.StepsStatus, error) {
		functionCalledStack = append(functionCalledStack, "findStepStatusByServiceRequestIdAndStatus")
		stepsStatus := make([]models.StepsStatus, 1)
		stepsStatus[0].StepName = "firstStep"
		stepsStatus[0].Payload.Request = map[string]interface{}{"k": "v"}
		stepsStatus[0].Payload.Response = map[string]interface{}{"k": "v"}
		return stepsStatus, nil
	}
	type args struct {
		serviceReq models.ServiceRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should continue execution from next step",
			args: args{
				serviceReq: models.ServiceRequest{
					ID:            uuid.New(),
					WorkflowName:  workflowName,
					Status:        models.STATUS_NEW,
					Payload:       prepareRequestPayload(),
					CurrentStepID: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddServiceRequestToChannel(tt.args.serviceReq)
			time.Sleep(time.Second * 5)
			assert.Equal(t, 0, len(serviceRequestChannel))
			assert.Equal(t, 4, len(functionCalledStack))
		})
	}
}
