package models

import (
	"clamp-core/executors"
	"clamp-core/transform"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func prepareStepRequestResponse() map[string]*StepContext {
	stepRequestResponse := map[string]*StepContext{"dummyStep": {
		Request:  map[string]interface{}{"user_type": "admin"},
		Response: nil,
	}}
	return stepRequestResponse
}

func prepareRequestContextForTests() RequestContext {
	reqCtx := RequestContext{
		ServiceRequestId: uuid.UUID{},
		WorkflowName:     "",
		StepsContext:     prepareStepRequestResponse(),
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
		requestBody    StepRequest
		prefix         string
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
				requestContext: RequestContext{
					ServiceRequestId: uuid.UUID{},
					WorkflowName:     "",
					StepsContext:     prepareStepRequestResponse(),
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
				Type:           tt.fields.StepType,
				Mode:           tt.fields.Mode,
				Val:            tt.fields.Val,
				Transform:      tt.fields.Transform,
				Enabled:        tt.fields.Enabled,
				When:           tt.fields.When,
				canStepExecute: tt.fields.canStepExecute,
			}

			reqCtx := prepareRequestContextForTests()
			_, err := step.DoExecute(reqCtx, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoExecute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.want(*step)
		})
	}
}

func TestStep_DoTransform(t *testing.T) {
	type fields struct {
		Id               int
		Name             string
		StepType         string
		Mode             string
		Val              Val
		Transform        bool
		Enabled          bool
		When             string
		canStepExecute   bool
		TransformFormat  string
		RequestTransform RequestTransform
	}
	type args struct {
		requestBody map[string]interface{}
		prefix      string
	}
	tests := []struct {
		name                   string
		fields                 fields
		args                   args
		expectedTransformation map[string]interface{}
		wantErr                bool
	}{
		{
			name: "ShouldTransformRequestWhenTransformIsEnabledAndJsonTransformation",
			fields: fields{
				Id:       1,
				Name:     "dummyStep",
				Mode:     "HTTP",
				StepType: "SYNC",
				RequestTransform: &transform.JsonTransform{
					Spec: map[string]interface{}{"userdetails.name": "user_name"},
				},
				Val: &executors.HttpVal{
					Method:  "POST",
					Url:     "https://reqres.in/api/users",
					Headers: "",
				},
				Transform: true,
				Enabled:   true,
			},
			args: args{
				requestBody: map[string]interface{}{"user_type": "admin", "user_name": "superadmin"},
				prefix:      "",
			},
			expectedTransformation: map[string]interface{}{"userdetails": map[string]interface{}{"name": "superadmin"}},
			wantErr:                false,
		},
		{
			name: "ShouldTransformRequestWhenTransformIsEnabledAndXmlTransformation",
			fields: fields{
				Id:       1,
				Name:     "dummyStep",
				Mode:     "HTTP",
				StepType: "SYNC",
				RequestTransform: &transform.XMLTransform{
					Keys: map[string]interface{}{"name": "user_name"},
				},
				TransformFormat: "XML",
				Val: &executors.HttpVal{
					Method:  "POST",
					Url:     "https://reqres.in/api/users",
					Headers: "",
				},
				Transform: true,
				Enabled:   true,
			},
			args: args{
				requestBody: map[string]interface{}{"user_type": "admin", "user_name": "superadmin"},
				prefix:      "",
			},
			expectedTransformation: map[string]interface{}{"user_type": "admin", "user_name": "superadmin"},
			wantErr:                false,
		},
		{
			name: "ShouldNotTransformRequestWhenTransformIsDisabled",
			fields: fields{
				Id:       1,
				Name:     "dummyStep",
				Mode:     "HTTP",
				StepType: "SYNC",
				RequestTransform: &transform.JsonTransform{
					Spec: map[string]interface{}{"name": "user_name"},
				},
				Val: &executors.HttpVal{
					Method:  "POST",
					Url:     "https://reqres.in/api/users",
					Headers: "",
				},
				Transform: false,
				Enabled:   true,
			},
			args: args{
				requestBody: map[string]interface{}{"user_type": "admin", "user_name": "superadmin"},
				prefix:      "",
			},
			expectedTransformation: map[string]interface{}{"user_type": "admin", "user_name": "superadmin"},
			wantErr:                false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			step := &Step{
				Id:               tt.fields.Id,
				Name:             tt.fields.Name,
				Type:             tt.fields.StepType,
				Mode:             tt.fields.Mode,
				Val:              tt.fields.Val,
				Transform:        tt.fields.Transform,
				TransformFormat:  tt.fields.TransformFormat,
				RequestTransform: tt.fields.RequestTransform,
				Enabled:          tt.fields.Enabled,
				When:             tt.fields.When,
				canStepExecute:   tt.fields.canStepExecute,
			}

			transformedRequest, err := step.DoTransform(tt.args.requestBody, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoTransform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			eq := reflect.DeepEqual(transformedRequest, tt.expectedTransformation)
			if !eq {
				t.Errorf("TransformRequest() transformedRequest = %v, want %v", transformedRequest, tt.expectedTransformation)
			}
		})
	}
}
