package services

import (
	"clamp-core/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const workflowName string = "testWF"

func setUp() {
	whereQueryMock = func(model interface{}, cond string, params ...interface{}) error {
		serviceReq := model.(*models.Workflow)
		serviceReq.Id = "TEST_WF"
		serviceReq.Steps = []models.Step{}
		return nil
	}
}

func TestAddServiceRequestToChannel(t *testing.T) {
	setUp()
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
		})
	}
}
