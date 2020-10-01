package hooks

import (
	"testing"
)

func setupStepRequest() map[string]interface{} {
	stepRequestResponsePayload := make(map[string]interface{})
	stepRequestResponsePayload["dummyStep"] = map[string]interface{}{"request": map[string]interface{}{"user_type": "admin"}, "response": map[string]interface{}{"user_type": "admin"}}
	return stepRequestResponsePayload
}

func TestExprHook_ShouldStepExecute(t *testing.T) {
	type args struct {
		whenCondition string
		stepRequest   map[string]interface{}
		prefix        string
	}
	tests := []struct {
		name               string
		args               args
		wantCanStepExecute bool
		wantErr            bool
	}{
		{
			name: "shouldReturnTrueIfConditionSatisfies",
			args: args{
				whenCondition: "dummyStep.request.user_type == 'admin'",
				stepRequest:   setupStepRequest(),
				prefix:        "",
			},
			wantCanStepExecute: true,
			wantErr:            false,
		}, {
			name: "shouldReturnFalseIfConditionNotSatisfied",
			args: args{
				whenCondition: "dummyStep.request.user_type == 'user'",
				stepRequest:   setupStepRequest(),
				prefix:        "",
			},
			wantCanStepExecute: false,
			wantErr:            false,
		}, {
			name: "shouldReturnErrorIfConditionIsNotInProperFormat",
			args: args{
				whenCondition: "user_type",
				stepRequest:   setupStepRequest(),
				prefix:        "",
			},
			wantCanStepExecute: false,
			wantErr:            true,
		}, {
			name: "shouldReturnErrorIfConditionIsInvalid",
			args: args{
				whenCondition: "1+2",
				stepRequest:   setupStepRequest(),
				prefix:        "",
			},
			wantCanStepExecute: false,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ExprHook{}
			gotCanStepExecute, err := e.ShouldStepExecute(tt.args.whenCondition, tt.args.stepRequest, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShouldStepExecute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCanStepExecute != tt.wantCanStepExecute {
				t.Errorf("ShouldStepExecute() gotCanStepExecute = %v, want %v", gotCanStepExecute, tt.wantCanStepExecute)
			}
		})
	}
}
