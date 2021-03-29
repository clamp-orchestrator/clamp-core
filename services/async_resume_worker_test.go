package services

import (
	"clamp-core/executors"
	"clamp-core/models"
	"clamp-core/utils"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var serviceRequestID = uuid.New()

func prepareRequestPayload() map[string]interface{} {
	var payload = map[string]interface{}{
		"id1": "val1",
		"id2": "val2",
	}
	return payload
}

func prepareResponsePayload() map[string]interface{} {
	var payload = map[string]interface{}{
		"key1": "val1",
		"key2": "val2",
	}
	return payload
}

func TestShouldAddSuccessResponseFromAsyncStepResponseToChannel(t *testing.T) {
	var functionCalledStack []string
	mockDB.FindWorkflowByNameMockFunc = func(workflowName string) (workflow *models.Workflow, err error) {
		workflow = &models.Workflow{}
		workflow.ID = "TEST_WF"
		step := models.Step{
			ID:        1,
			Name:      "1",
			Type:      utils.StepTypeAsync,
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
		return stepStatus, nil
	}
	type args struct {
		asyncStepResponseReq models.AsyncStepResponse
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should process Async Step Response request",
			args: args{
				asyncStepResponseReq: models.AsyncStepResponse{
					ServiceRequestID: serviceRequestID,
					StepID:           1,
					Response:         prepareResponsePayload(),
					Error:            models.ClampErrorResponse{},
				},
			},
		},
	}

	mockDB.FindServiceRequestByIDMockFunc = func(u uuid.UUID) (request *models.ServiceRequest, err error) {
		serviceRequest := &models.ServiceRequest{
			ID:            serviceRequestID,
			WorkflowName:  workflowName,
			Status:        models.StatusNew,
			CreatedAt:     time.Time{},
			Payload:       prepareRequestPayload(),
			CurrentStepID: 1,
		}
		functionCalledStack = append(functionCalledStack, "findServiceRequestById")
		return serviceRequest, err
	}

	mockDB.FindAllStepStatusByServiceRequestIDAndStepIDMockFunc = func(serviceRequestId uuid.UUID, stepId int) (stepsStatus []*models.StepsStatus, err error) {
		var statuses = make([]*models.StepsStatus, 1)
		stepStatus := models.StepsStatus{
			ID:               "2",
			ServiceRequestID: serviceRequestId,
			WorkflowName:     workflowName,
			Status:           models.StatusStarted,
			CreatedAt:        time.Now(),
			TotalTimeInMs:    1000,
			StepName:         "1",
			Reason:           "",
			Payload: models.Payload{
				Request:  prepareRequestPayload(),
				Response: nil,
			},
			StepID: 1,
		}
		statuses[0] = &stepStatus
		functionCalledStack = append(functionCalledStack, "findStepStatusByServiceRequestIdAndStepIdAndStatus")
		return statuses, err
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddStepResponseToResumeChannel(&tt.args.asyncStepResponseReq)
			time.Sleep(time.Millisecond)
			assert.Equal(t, 0, len(serviceRequestChannel))
			assert.Equal(t, 4, len(functionCalledStack))
		})
	}
}

func TestShouldAddFailureResponseFromAsyncStepResponseToChannel(t *testing.T) {
	var functionCalledStack []string
	mockDB.FindWorkflowByNameMockFunc = func(workflowName string) (workflow *models.Workflow, err error) {
		workflow = &models.Workflow{}
		workflow.ID = "TEST_WF"
		step := models.Step{
			ID:        1,
			Name:      "1",
			Type:      utils.StepTypeAsync,
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
		return stepStatus, nil
	}
	type args struct {
		asyncStepResponseReq models.AsyncStepResponse
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should process Async Step Response request",
			args: args{
				asyncStepResponseReq: models.AsyncStepResponse{
					ServiceRequestID: serviceRequestID,
					StepID:           1,
					Response:         prepareResponsePayload(),
					Error: models.ClampErrorResponse{
						Code:    http.StatusBadRequest,
						Message: "Failed to process due to internal failure",
					},
				},
			},
		},
	}

	mockDB.FindAllStepStatusByServiceRequestIDAndStepIDMockFunc = func(serviceRequestId uuid.UUID, stepId int) (stepsStatus []*models.StepsStatus, err error) {
		var statuses = make([]*models.StepsStatus, 1)
		stepStatus := models.StepsStatus{
			ID:               "2",
			ServiceRequestID: serviceRequestId,
			WorkflowName:     workflowName,
			Status:           models.StatusStarted,
			CreatedAt:        time.Now(),
			TotalTimeInMs:    1000,
			StepName:         "1",
			Reason:           "",
			Payload: models.Payload{
				Request:  prepareRequestPayload(),
				Response: nil,
			},
			StepID: 1,
		}
		statuses[0] = &stepStatus
		functionCalledStack = append(functionCalledStack, "findStepStatusByServiceRequestIdAndStepIdAndStatus")
		return statuses, err
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddStepResponseToResumeChannel(&tt.args.asyncStepResponseReq)
			time.Sleep(time.Millisecond)
			assert.Equal(t, 0, len(serviceRequestChannel))
			assert.Equal(t, 2, len(functionCalledStack))
		})
	}
}
