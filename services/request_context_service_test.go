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
	findStepStatusByServiceRequestIdAndStatusMock = func(serviceRequestId uuid.UUID, status models.Status) ([]models.StepsStatus, error) {
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

func TestComputeRequestToCurrentStepInContext(t *testing.T) {
	workflow := models.Workflow{
		Name: "TEST_WF",
		Steps: []models.Step{{
			Name: "step1",
		}, {
			Name: "step2",
		}, {
			Name: "step3",
		}},
	}
	context := CreateRequestContext(workflow, models.ServiceRequest{
		ID: uuid.New(),
	})

	type args struct {
		workflow             models.Workflow
		currentStepExecuting models.Step
		requestContext       *models.RequestContext
		stepIndex            int
		stepRequestPayload   map[string]interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ShouldComputeContextForFirstStep",
			args: args{
				workflow: workflow,
				currentStepExecuting: models.Step{
					Name: "step1",
				},
				requestContext:     &context,
				stepIndex:          0,
				stepRequestPayload: map[string]interface{}{"k": "v"},
			},
		},
		{
			name: "ShouldComputeContextForIntermediateStep2",
			args: args{
				workflow: workflow,
				currentStepExecuting: models.Step{
					Name: "step2",
				},
				requestContext:     &context,
				stepIndex:          1,
				stepRequestPayload: map[string]interface{}{"request": "value"},
			},
		},
		{
			name: "ShouldComputeContextForIntermediateStep3",
			args: args{
				workflow: workflow,
				currentStepExecuting: models.Step{
					Name: "step3",
				},
				requestContext:     &context,
				stepIndex:          1,
				stepRequestPayload: map[string]interface{}{"request": "value"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ComputeRequestToCurrentStepInContext(tt.args.workflow, tt.args.currentStepExecuting, tt.args.requestContext, tt.args.stepIndex, tt.args.stepRequestPayload)
			tt.args.requestContext.SetStepResponseToContext(tt.args.currentStepExecuting.Name, map[string]interface{}{"response": "value"})
			assert.NotNil(t, tt.args.requestContext.GetStepRequestFromContext(tt.args.currentStepExecuting.Name))
		})
	}
}
