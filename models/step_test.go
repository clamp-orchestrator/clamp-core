package models

import (
	"clamp-core/executors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func prepareStepRequestResponse() map[string]RequestResponse {
	stepRequestResponse := map[string]RequestResponse{"dummyStep": {
		Request:  map[string]interface{}{"user_type": "admin"},
		Response: nil,
	}}
	return stepRequestResponse
}

func prepareRequestContextForTests() RequestContext {
	reqCtx := RequestContext{
		ServiceRequestId: uuid.UUID{},
		WorkflowName:     "",
		Payload:          prepareStepRequestResponse(),
	}
	return reqCtx
}

func TestStep_DoExecute(t *testing.T) {
	type fields struct {
		Id             int
		Name           string
		StepType       string
		Mode           string
		Val            Val
		Transform      bool
		Enabled        bool
		When           string
		canStepExecute bool
	}
	type args struct {
		requestBody StepRequest
		prefix      string
		requestContext RequestContext
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    func(step Step)
		wantErr bool
	}{
		{
			name: "ShouldExecuteStepIfWhenConditionSatisfied",
			fields: fields{
				Id:       1,
				Name:     "dummyStep",
				Mode:     "HTTP",
				StepType: "SYNC",
				Val: &executors.HttpVal{
					Method: "POST",
					Url:    "http://34.222.238.234:3333/api/v1/orders",
				},
				Transform: false,
				Enabled:   true,
				When:      "context.dummyStep.request.user_type == 'admin'",
			},
			args: args{
				requestBody: StepRequest{
					ServiceRequestId: uuid.UUID{},
					StepId:           0,
					Payload:          map[string]interface{}{"user_type": "admin"},
				},
				prefix: "",
			},
			want: func(step Step) {
				assert.True(t, step.canStepExecute)
			},
			wantErr: false,
		},
		{
			name: "ShouldNotExecuteStepIfWhenConditionNotSatisfied",
			fields: fields{
				Id:       1,
				Name:     "dummyStep",
				Mode:     "HTTP",
				StepType: "SYNC",
				Val: &executors.HttpVal{
					Method: "POST",
					Url:    "http://34.222.238.234:3333/api/v1/orders",
				},
				Transform: false,
				Enabled:   true,
				When:      "context.dummyStep.request.user_type == 'user'",
			},
			args: args{
				requestBody: StepRequest{
					ServiceRequestId: uuid.UUID{},
					StepId:           0,
					Payload:          map[string]interface{}{"user_type": "admin"},
				},
				prefix: "",
				requestContext : RequestContext{
					ServiceRequestId: uuid.UUID{},
					WorkflowName:     "",
					Payload:          prepareStepRequestResponse(),
				},
			},
			want: func(step Step) {
				assert.False(t, step.canStepExecute)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			step := &Step{
				Id:             tt.fields.Id,
				Name:           tt.fields.Name,
				StepType:       tt.fields.StepType,
				Mode:           tt.fields.Mode,
				Val:            tt.fields.Val,
				Transform:      tt.fields.Transform,
				Enabled:        tt.fields.Enabled,
				When:           tt.fields.When,
				canStepExecute: tt.fields.canStepExecute,
			}

			reqCtx := prepareRequestContextForTests()
			_, err := step.DoExecute(tt.args.requestBody, tt.args.prefix, reqCtx)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoExecute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.want(*step)
		})
	}
}
