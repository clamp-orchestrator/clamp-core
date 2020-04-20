package services

import (
	"clamp-core/executors"
	"clamp-core/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var step models.Step

func TestAddAsyncStepExecutorRequestToChannel(t *testing.T) {
	var functionCalledStack []string
	findWorkflowByNameMock = func(workflowName string) (workflow models.Workflow, err error) {
		workflow.Id = "TEST_WF"
		step := models.Step{
			Id:        1,
			Name:      "1",
			StepType:  "ASYNC",
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
	findAllStepStatusByServiceRequestIdAndStepIdMock = func(serviceRequestId uuid.UUID, stepId int) (statuses []models.StepsStatus, err error) {
		return statuses, err
	}
	type args struct {
		serviceReq models.ServiceRequest
	}
	var serviceRequestId = uuid.New()
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should process service request",
			args: args{
				serviceReq: models.ServiceRequest{
					ID:           serviceRequestId,
					WorkflowName: workflowName,
					Status:       models.STATUS_NEW,
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
			Payload:       nil,
			CurrentStepId: 1,
		}
		functionCalledStack = append(functionCalledStack, "findServiceRequestById")
		return serviceRequest, err
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddServiceRequestToChannel(tt.args.serviceReq)
			time.Sleep(time.Second * 5)
			assert.Equal(t, 0, len(serviceRequestChannel))
			assert.Equal(t, 5, len(functionCalledStack))
		})
	}
}
