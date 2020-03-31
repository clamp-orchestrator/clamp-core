package services

import (
	"clamp-core/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAddServiceRequestToChannel1(t *testing.T) {
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
					WorkflowName: "TESTING",
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
