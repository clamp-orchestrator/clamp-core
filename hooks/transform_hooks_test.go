package hooks

import (
	"reflect"
	"testing"
)

func TestExprHook_TransformRequest(t *testing.T) {
	type args struct {
		key         map[string]interface{}
		stepRequest map[string]interface{}
	}
	tests := []struct {
		name                   string
		args                   args
		expectedTransformation map[string]interface{}
		wantErr                bool
		errMessage             string
	}{
		{
			name: "shouldReturnTransformedValueIfKeysMatchesWithRequestPayload",
			args: args{
				key:         map[string]interface{}{"userType": "dummyStep.request.user_type"},
				stepRequest: setupStepRequest(),
			},
			expectedTransformation: map[string]interface{}{"userType": "admin"},
			wantErr:                false,
		},
		{
			name: "shouldReturnTransformationFailedWhenSpecIsEmpty",
			args: args{
				key:         map[string]interface{}{},
				stepRequest: setupStepRequest(),
			},
			wantErr:    true,
			errMessage: "SpecError - Spec must contain at least one element",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &TransformHook{}
			transformedRequest, err := e.TransformRequest(tt.args.stepRequest, tt.args.key)

			if (err != nil) != tt.wantErr {
				t.Errorf("TransformRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				eq := reflect.DeepEqual(transformedRequest, tt.expectedTransformation)
				if !eq {
					t.Errorf("TransformRequest() transformedRequest = %v, want %v", transformedRequest, tt.expectedTransformation)
				}
			}
		})
	}
}
