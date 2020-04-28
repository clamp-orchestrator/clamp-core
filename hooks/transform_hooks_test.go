package hooks

import (
	"testing"
)

func TestExprHook_TransformRequest(t *testing.T) {
	type args struct {
		key string
		stepRequest   map[string]interface{}
	}
	tests := []struct {
		name             string
		args             args
		transformedValue string
		wantErr          bool
	}{
		{
			name: "shouldReturnTransformedValueIfKeysMatchesWithRequestPayload",
			args: args{
				key: "dummyStep.request.user_type",
				stepRequest: setupStepRequest(),
			},
			transformedValue: "admin",
			wantErr:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &TransformHook{}
			transformedRequest, err := e.TransformRequest( tt.args.stepRequest, tt.args.key)
			transformedRequestValue := transformedRequest[tt.args.key].(string)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransformRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if transformedRequestValue != tt.transformedValue {
				t.Errorf("TransformRequest() transformedRequestValue = %v, want %v", transformedRequestValue, tt.transformedValue)
			}
		})
	}
}
