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

type mockRepoImpl struct {
}

func (s mockRepoImpl) whereQuery(model interface{}, condition string, params ...interface{}) error {
	serviceReq := model.(*models.Workflow)
	serviceReq.Id = "TEST_WF"
	step := models.Step{
		Id:        "1",
		Name:      "1",
		Mode:      "HTTP",
		Transform: false,
		Enabled:   false,
		Val: &executors.HttpVal{
			Method:  "POST",
			Url:     "http://35.166.176.234:3333/api/v1/login",
			Headers: "",
		},
	}
	serviceReq.Steps = []models.Step{step}

	return nil
}

func (s mockRepoImpl) insertQuery(model interface{}) error {
	return nil
}

func (s mockRepoImpl) selectQuery(model interface{}) error {
	return nil
}

func (s mockRepoImpl) query(model interface{}, query interface{}, param interface{}) (Result, error) {
	panic("Query func not implemented")
}

func setUp() {
	repo = mockRepoImpl{}
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
