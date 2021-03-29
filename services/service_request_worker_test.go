package services

import (
	"clamp-core/executors"
	"clamp-core/models"
	"clamp-core/transform"
	"clamp-core/utils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const workflowName string = "testWF"

func TestAddServiceRequestToChannel(t *testing.T) {
	var functionCalledStack []string
	mockDB.FindWorkflowByNameMockFunc = func(workflowName string) (workflow *models.Workflow, err error) {
		workflow = &models.Workflow{}
		workflow.ID = "TEST_WF"
		step := models.Step{
			Name:      "1",
			Type:      utils.StepTypeSync,
			Mode:      utils.StepModeHTTP,
			Transform: false,
			Enabled:   false,
			Val: &executors.HTTPVal{
				Method:  "POST",
				URL:     "https://run.mocky.io/v3/0590fbf8-0f1c-401c-b9df-65e98ef0385d",
				Headers: "",
			},
		}
		workflow.Steps = []models.Step{step}
		functionCalledStack = append(functionCalledStack, "findWorkflowByName")
		return workflow, err
	}
	mockDB.SaveStepStatusMockFunc = func(stepStatus *models.StepsStatus) (status *models.StepsStatus, err error) {
		functionCalledStack = append(functionCalledStack, "saveStepStatusMock")
		status = &models.StepsStatus{}
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
					Status:       models.StatusNew,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddServiceRequestToChannel(&tt.args.serviceReq)
			time.Sleep(time.Second * 5)
			assert.Equal(t, 0, len(serviceRequestChannel))
			assert.Equal(t, 3, len(functionCalledStack))
		})
	}
}

func TestShouldAddServiceRequestToChannelWithTransformationEnabledForOneStepInTheWorkflow(t *testing.T) {
	var functionCalledStack []string
	mockDB.FindWorkflowByNameMockFunc = func(workflowName string) (workflow *models.Workflow, err error) {
		workflow = &models.Workflow{}
		workflow.ID = "TEST_WF"
		step := models.Step{
			Name:      "1",
			Type:      utils.StepTypeSync,
			Mode:      utils.StepModeHTTP,
			Transform: true,
			Enabled:   false,
			RequestTransform: &transform.JSONTransform{
				Spec: map[string]interface{}{"name": "test"},
			},
			Val: &executors.HTTPVal{
				Method:  "POST",
				URL:     "https://run.mocky.io/v3/0590fbf8-0f1c-401c-b9df-65e98ef0385d",
				Headers: "",
			},
		}
		workflow.Steps = []models.Step{step}
		functionCalledStack = append(functionCalledStack, "findWorkflowByName")
		return workflow, err
	}
	mockDB.SaveStepStatusMockFunc = func(stepStatus *models.StepsStatus) (status *models.StepsStatus, err error) {
		functionCalledStack = append(functionCalledStack, "saveStepStatusMock")
		status = &models.StepsStatus{}
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
					Status:       models.StatusNew,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddServiceRequestToChannel(&tt.args.serviceReq)
			time.Sleep(time.Second * 5)
			assert.Equal(t, 0, len(serviceRequestChannel))
			assert.Equal(t, 3, len(functionCalledStack))
		})
	}
}

func TestShouldSkipStepIfConditionDoesNotMatch(t *testing.T) {
	var functionCalledStack []string
	mockDB.FindWorkflowByNameMockFunc = func(workflowName string) (workflow *models.Workflow, err error) {
		workflow = &models.Workflow{}
		workflow.ID = "TEST_WF"
		step := models.Step{
			Name:      "skipStep",
			Type:      utils.StepTypeSync,
			Mode:      utils.StepModeHTTP,
			Transform: false,
			Enabled:   false,
			When:      "skipStep.request.id1 == 'val3'",
			Val: &executors.HTTPVal{
				Method:  "POST",
				URL:     "https://run.mocky.io/v3/0590fbf8-0f1c-401c-b9df-65e98ef0385d",
				Headers: "",
			},
		}
		workflow.Steps = []models.Step{step}
		functionCalledStack = append(functionCalledStack, "findWorkflowByName")
		return workflow, err
	}
	mockDB.SaveStepStatusMockFunc = func(stepStatus *models.StepsStatus) (status *models.StepsStatus, err error) {
		functionCalledStack = append(functionCalledStack, "saveStepStatusMock")
		status = &models.StepsStatus{}
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
					Status:       models.StatusNew,
					Payload:      prepareRequestPayload(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddServiceRequestToChannel(&tt.args.serviceReq)
			time.Sleep(time.Second * 5)
			assert.Equal(t, 0, len(serviceRequestChannel))
			assert.Equal(t, 3, len(functionCalledStack))
		})
	}
}

func TestShouldResumeTheWorkflowExecutionFromNextStep(t *testing.T) {
	var functionCalledStack []string
	mockDB.FindWorkflowByNameMockFunc = func(workflowName string) (workflow *models.Workflow, err error) {
		workflow = &models.Workflow{}
		workflow.ID = "TEST_WF"
		step := models.Step{
			Name:      "firstStep",
			Type:      utils.StepTypeSync,
			Mode:      utils.StepModeHTTP,
			Transform: false,
			Enabled:   false,
			When:      "firstStep.request.id1 == 'val1'",
			Val: &executors.HTTPVal{
				Method:  "POST",
				URL:     "https://run.mocky.io/v3/0590fbf8-0f1c-401c-b9df-65e98ef0385d",
				Headers: "",
			},
		}
		step1 := models.Step{
			Name:      "secondStep",
			Type:      utils.StepTypeSync,
			Mode:      utils.StepModeHTTP,
			Transform: false,
			Enabled:   false,
			When:      "firstStep.request.id1 == 'val1'",
			Val: &executors.HTTPVal{
				Method:  "POST",
				URL:     "https://run.mocky.io/v3/0590fbf8-0f1c-401c-b9df-65e98ef0385d",
				Headers: "",
			},
		}
		workflow.Steps = []models.Step{step, step1}
		functionCalledStack = append(functionCalledStack, "findWorkflowByName")
		return workflow, err
	}
	mockDB.SaveStepStatusMockFunc = func(stepStatus *models.StepsStatus) (status *models.StepsStatus, err error) {
		functionCalledStack = append(functionCalledStack, "saveStepStatusMock")
		status = &models.StepsStatus{}
		return status, nil
	}
	mockDB.FindStepStatusByServiceRequestIDAndStatusMockFunc = func(serviceRequestId uuid.UUID, status models.Status) ([]*models.StepsStatus, error) {
		functionCalledStack = append(functionCalledStack, "findStepStatusByServiceRequestIdAndStatus")
		stepsStatus := make([]*models.StepsStatus, 1)
		stepsStatus[0] = &models.StepsStatus{}
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
					Status:        models.StatusNew,
					Payload:       prepareRequestPayload(),
					CurrentStepID: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddServiceRequestToChannel(&tt.args.serviceReq)
			time.Sleep(time.Second * 5)
			assert.Equal(t, 0, len(serviceRequestChannel))
			assert.Equal(t, 4, len(functionCalledStack))
		})
	}
}
