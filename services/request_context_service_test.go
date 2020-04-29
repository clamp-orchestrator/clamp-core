package services

import (
	"clamp-core/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCreateRequestContext(t *testing.T) {
	type args struct {
		workflow models.Workflow
		request  models.ServiceRequest
	}
	id := uuid.New()
	tests := []struct {
		name        string
		args        args
		wantContext models.RequestContext
	}{
		{
			name: "ShouldCreateRequestContext",
			args: args{
				workflow: models.Workflow{
					Name:  "TEST_WF",
					Steps: []models.Step{{Name: "step1"}},
				},
				request: models.ServiceRequest{
					ID: id,
				},
			},
			wantContext: models.RequestContext{
				ServiceRequestId: id,
				WorkflowName:     "TEST_WF",
				StepsContext: map[string]*models.StepContext{"step1": {
					Request:  nil,
					Response: nil,
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotContext := CreateRequestContext(tt.args.workflow, tt.args.request); !reflect.DeepEqual(gotContext, tt.wantContext) {
				t.Errorf("CreateRequestContext() = %v, want %v", gotContext, tt.wantContext)
			}
		})
	}
}

func TestShouldEnhanceRequestContext(t *testing.T) {
	context := CreateRequestContext(models.Workflow{
		Name: "TEST_WF",
		Steps: []models.Step{{
			Name: "step1",
		}, {
			Name: "step2",
		}},
	}, models.ServiceRequest{
		ID: uuid.New(),
	})
	findStepStatusByServiceRequestIdAndStatusOrderByCreatedAtDescMock = func(serviceRequestId uuid.UUID, status models.Status) ([]models.StepsStatus, error) {
		stepsStatus := make([]models.StepsStatus, 2)
		stepsStatus[0].StepName = "step1"
		stepsStatus[0].Payload.Request = map[string]interface{}{"k": "v"}
		stepsStatus[0].Payload.Response = map[string]interface{}{"k": "v"}
		stepsStatus[1].StepName = "step2"
		stepsStatus[1].Payload.Request = map[string]interface{}{"k": "v"}
		stepsStatus[1].Payload.Response = map[string]interface{}{"k": "v"}
		return stepsStatus, nil
	}

	EnhanceRequestContextWithExecutedSteps(&context)
	if ctx, ok := context.StepsContext["step1"]; true {
		assert.True(t, ok)
		assert.NotNil(t, ctx.Request)
		assert.NotNil(t, ctx.Response)
	}
	if ctx, ok := context.StepsContext["step2"]; true {
		assert.True(t, ok)
		assert.NotNil(t, ctx.Request)
		assert.NotNil(t, ctx.Response)
	}

}
