package services

import (
	"clamp-core/executors"
	"clamp-core/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var serviceRequestId = uuid.New()

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
	findWorkflowByNameMock = func(workflowName string) (workflow models.Workflow, err error) {
		workflow.Id = "TEST_WF"
		step := models.Step{
			Id:        1,
			Name:      "1",
			Type:      "ASYNC",
			Mode:      "HTTP",
			Transform: false,
			Enabled:   false,
			Val: &executors.HttpVal{
				Method:  "POST",
				Url:     "http://34.216.32.148:3333/api/v1/login",
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
					ServiceRequestId: serviceRequestId,
					StepId:           1,
					Response:         prepareResponsePayload(),
					Error:            models.ClampErrorResponse{},
				},
			},
		},
	}

	findServiceRequestByIdMock = func(u uuid.UUID) (request models.ServiceRequest, err error) {
		serviceRequest := models.ServiceRequest{
			ID:            serviceRequestId,
			WorkflowName:  workflowName,
			Status:        models.STATUS_NEW,
			CreatedAt:     time.Time{},
			Payload:       prepareRequestPayload(),
			CurrentStepId: 1,
		}
		functionCalledStack = append(functionCalledStack, "findServiceRequestById")
		return serviceRequest, err
	}

	findAllStepStatusByServiceRequestIdAndStepIdMock = func(serviceRequestId uuid.UUID, stepId int) (stepsStatus []models.StepsStatus, err error) {
		var statuses = make([]models.StepsStatus, 1)
		stepStatus := models.StepsStatus{
			ID:               "2",
			ServiceRequestId: serviceRequestId,
			WorkflowName:     workflowName,
			Status:           models.STATUS_STARTED,
			CreatedAt:        time.Now(),
			TotalTimeInMs:    1000,
			StepName:         "1",
			Reason:           "",
			Payload: models.Payload{
				Request:  prepareRequestPayload(),
				Response: nil,
			},
			StepId: 1,
		}
		statuses[0] = stepStatus
		functionCalledStack = append(functionCalledStack, "findStepStatusByServiceRequestIdAndStepIdAndStatus")
		return statuses, err
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddStepResponseToResumeChannel(tt.args.asyncStepResponseReq)
			time.Sleep(time.Second * 5)
			assert.Equal(t, 0, len(serviceRequestChannel))
			assert.Equal(t, 4, len(functionCalledStack))
		})
	}
}

func TestShouldAddFailureResponseFromAsyncStepResponseToChannel(t *testing.T) {
	var functionCalledStack []string
	findWorkflowByNameMock = func(workflowName string) (workflow models.Workflow, err error) {
		workflow.Id = "TEST_WF"
		step := models.Step{
			Id:        1,
			Name:      "1",
			Type:      "ASYNC",
			Mode:      "HTTP",
			Transform: false,
			Enabled:   false,
			Val: &executors.HttpVal{
				Method:  "POST",
				Url:     "http://34.216.32.148:3333/api/v1/login",
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
					ServiceRequestId: serviceRequestId,
					StepId:           1,
					Response:         prepareResponsePayload(),
					Error: models.ClampErrorResponse{
						Code:    400,
						Message: "Failed to process due to internal failure",
					},
				},
			},
		},
	}

	findAllStepStatusByServiceRequestIdAndStepIdMock = func(serviceRequestId uuid.UUID, stepId int) (stepsStatus []models.StepsStatus, err error) {
		var statuses = make([]models.StepsStatus, 1)
		stepStatus := models.StepsStatus{
			ID:               "2",
			ServiceRequestId: serviceRequestId,
			WorkflowName:     workflowName,
			Status:           models.STATUS_STARTED,
			CreatedAt:        time.Now(),
			TotalTimeInMs:    1000,
			StepName:         "1",
			Reason:           "",
			Payload: models.Payload{
				Request:  prepareRequestPayload(),
				Response: nil,
			},
			StepId: 1,
		}
		statuses[0] = stepStatus
		functionCalledStack = append(functionCalledStack, "findStepStatusByServiceRequestIdAndStepIdAndStatus")
		return statuses, err
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddStepResponseToResumeChannel(tt.args.asyncStepResponseReq)
			time.Sleep(time.Second * 5)
			assert.Equal(t, 0, len(serviceRequestChannel))
			assert.Equal(t, 2, len(functionCalledStack))
		})
	}
}
