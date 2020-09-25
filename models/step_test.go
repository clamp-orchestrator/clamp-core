package models

import (
	"clamp-core/executors"
	"clamp-core/transform"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func prepareStepRequestResponse() map[string]*StepContext {
	stepRequestResponse := map[string]*StepContext{"dummyStep": {
		Request:     map[string]interface{}{"user_type": "admin"},
		Response:    nil,
		StepSkipped: false,
	}}
	return stepRequestResponse
}

func prepareRequestContextForTests() RequestContext {
	reqCtx := RequestContext{
		ServiceRequestID: uuid.UUID{},
		WorkflowName:     "",
		StepsContext:     prepareStepRequestResponse(),
	}
	return reqCtx
}

func TestStep_DoExecute(t *testing.T) {
	type fields struct {
		ID             int
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
				ID:       1,
				Name:     "dummyStep",
				Mode:     "HTTP",
				StepType: "SYNC",
				Val: &executors.HTTPVal{
					Method: "POST",
					URL:    "https://run.mocky.io/v3/0590fbf8-0f1c-401c-b9df-65e98ef0385d",
				},
				Transform: false,
				Enabled:   true,
				When:      "context.dummyStep.request.user_type == 'admin'",
			},
			args: args{
				requestBody: StepRequest{
					ServiceRequestID: uuid.UUID{},
					StepID:           0,
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
				ID:       1,
				Name:     "dummyStep",
				Mode:     "HTTP",
				StepType: "SYNC",
				Val: &executors.HTTPVal{
					Method: "POST",
					URL:    "https://run.mocky.io/v3/0590fbf8-0f1c-401c-b9df-65e98ef0385d",
				},
				Transform: false,
				Enabled:   true,
				When:      "context.dummyStep.request.user_type == 'user'",
			},
			args: args{
				requestBody: StepRequest{
					ServiceRequestID: uuid.UUID{},
					StepID:           0,
					Payload:          map[string]interface{}{"user_type": "admin"},
				},
				prefix: "",
				requestContext: RequestContext{
					ServiceRequestID: uuid.UUID{},
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
				ID:             tt.fields.ID,
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
		ID               int
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
		reqCtx      RequestContext
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
				ID:       1,
				Name:     "dummyStep",
				Mode:     "HTTP",
				StepType: "SYNC",
				RequestTransform: &transform.JSONTransform{
					Spec: map[string]interface{}{"userdetails.name": "dummyStep.request.user_name"},
				},
				Val: &executors.HTTPVal{
					Method:  "POST",
					URL:     "https://reqres.in/api/users",
					Headers: "",
				},
				Transform: true,
				Enabled:   true,
			},
			args: args{
				reqCtx: RequestContext{
					ServiceRequestID: uuid.UUID{},
					WorkflowName:     "",
					StepsContext: map[string]*StepContext{"dummyStep": {
						Request:  map[string]interface{}{"user_type": "admin", "user_name": "superadmin"},
						Response: nil,
					}},
				},
				requestBody: map[string]interface{}{"user_type": "admin", "user_name": "superadmin"},
				prefix:      "",
			},
			expectedTransformation: map[string]interface{}{"userdetails": map[string]interface{}{"name": "superadmin"}},
			wantErr:                false,
		},
		{
			name: "ShouldTransformRequestWhenTransformIsEnabledAndXmlTransformation",
			fields: fields{
				ID:       1,
				Name:     "dummyStep",
				Mode:     "HTTP",
				StepType: "SYNC",
				RequestTransform: &transform.XMLTransform{
					Keys: map[string]interface{}{"name": "dummyStep.request.user_name"},
				},
				TransformFormat: "XML",
				Val: &executors.HTTPVal{
					Method:  "POST",
					URL:     "https://reqres.in/api/users",
					Headers: "",
				},
				Transform: true,
				Enabled:   true,
			},
			args: args{
				reqCtx: RequestContext{
					ServiceRequestID: uuid.UUID{},
					WorkflowName:     "",
					StepsContext: map[string]*StepContext{"dummyStep": {
						Request:  map[string]interface{}{"user_type": "admin", "user_name": "superadmin"},
						Response: map[string]interface{}{},
					}},
				},
				requestBody: map[string]interface{}{"user_type": "admin", "user_name": "superadmin"},
				prefix:      "",
			},
			expectedTransformation: map[string]interface{}{
				"dummyStep": map[string]interface{}{
					"request":  map[string]interface{}{"user_type": "admin", "user_name": "superadmin"},
					"response": map[string]interface{}{},
				},
			},
			wantErr: false,
		},
		{
			name: "ShouldNotTransformRequestWhenTransformIsDisabled",
			fields: fields{
				ID:       1,
				Name:     "dummyStep",
				Mode:     "HTTP",
				StepType: "SYNC",
				RequestTransform: &transform.JSONTransform{
					Spec: map[string]interface{}{"name": "dummyStep.request.user_name"},
				},
				Val: &executors.HTTPVal{
					Method:  "POST",
					URL:     "https://reqres.in/api/users",
					Headers: "",
				},
				Transform: false,
				Enabled:   true,
			},
			args: args{
				reqCtx: RequestContext{
					ServiceRequestID: uuid.UUID{},
					WorkflowName:     "",
					StepsContext: map[string]*StepContext{"dummyStep": {
						Request:  map[string]interface{}{"user_type": "admin", "user_name": "superadmin"},
						Response: map[string]interface{}{},
					}},
				},
				prefix: "",
			},
			expectedTransformation: map[string]interface{}{
				"dummyStep": map[string]interface{}{
					"request":  map[string]interface{}{"user_type": "admin", "user_name": "superadmin"},
					"response": map[string]interface{}{},
				},
			},
			wantErr: false,
		},
		{
			name: "ShouldNotTransformRequestWhenTransformationFailsDueToSomeError",
			fields: fields{
				ID:       1,
				Name:     "dummyStep",
				Mode:     "HTTP",
				StepType: "SYNC",
				RequestTransform: &transform.JSONTransform{
					Spec: map[string]interface{}{},
				},
				Val: &executors.HTTPVal{
					Method:  "POST",
					URL:     "https://reqres.in/api/users",
					Headers: "",
				},
				Transform: true,
				Enabled:   true,
			},
			args: args{
				reqCtx: RequestContext{
					ServiceRequestID: uuid.UUID{},
					WorkflowName:     "",
					StepsContext: map[string]*StepContext{"dummyStep": {
						Request:  map[string]interface{}{"user_type": "admin", "user_name": "superadmin"},
						Response: map[string]interface{}{},
					}},
				},
				requestBody: map[string]interface{}{"user_type": "admin", "user_name": "superadmin"},
				prefix:      "",
			},
			expectedTransformation: nil,
			wantErr:                true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			step := &Step{
				ID:               tt.fields.ID,
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
			transformedRequest, err := step.DoTransform(tt.args.reqCtx, tt.args.prefix)
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
