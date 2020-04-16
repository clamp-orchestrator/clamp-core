package services

import (
	"clamp-core/executors"
	"clamp-core/models"
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
			Id:        "1",
			Name:      "1",
			StepType:  "SYNC",
			Mode:      "HTTP",
			Transform: false,
			Enabled:   false,
			Val: &executors.HttpVal{
				Method:  "POST",
				Url:     "http://35.166.176.234:3333/api/v1/login",
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
			//TODO @Tejash i have changed to 2, before it was 3 functionCalledStack is called twice na
			assert.Equal(t, 2, len(functionCalledStack))
		})
	}
}
