package executors

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpVal_DoExecute(t *testing.T) {
	type args struct {
		requestBody interface{}
	}
	tests := []struct {
		name       string
		fields     HTTPVal
		args       args
		want       interface{}
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "TestShouldExecuteHTTPStep",
			fields: HTTPVal{
				Method:  "GET",
				URL:     "https://run.mocky.io/v3/0590fbf8-0f1c-401c-b9df-65e98ef0385d",
				Headers: "",
			},
			args: args{
				requestBody: map[string]interface{}{"k": "v"},
			},
			want:       map[string]interface{}{"id": "1234", "name": "ABC", "email": "abc@sahaj.com", "org": "sahaj"},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "TestShouldExecuteHTTPStepWithHeaders",
			fields: HTTPVal{
				Method:  "GET",
				URL:     "https://run.mocky.io/v3/0590fbf8-0f1c-401c-b9df-65e98ef0385d",
				Headers: "Content-Type:application/json;token:abc",
			},
			args: args{
				requestBody: map[string]interface{}{"k": "v"},
			},
			want:       map[string]interface{}{"id": "1234", "name": "ABC", "email": "abc@sahaj.com", "org": "sahaj"},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "TestShouldThrowErrorWhileExecutingStep",
			fields: HTTPVal{
				Method:  "GET",
				URL:     "https://run.mocky.io/v3/nonexistent",
				Headers: "",
			},
			args: args{
				requestBody: map[string]interface{}{"k": "v"},
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "Not found",
		},
		{
			name: "TestShouldThrowErrorForHTTPStep",
			fields: HTTPVal{
				Method:  "GET",
				URL:     "http://localhost:3333/api/v1/user",
				Headers: "",
			},
			args: args{
				requestBody: map[string]interface{}{"k": "v"},
			},
			want:       map[string]interface{}{"id": "1234", "name": "ABC", "email": "abc@sahaj.com", "org": "sahaj"},
			wantErr:    true,
			wantErrMsg: "connect: connection refused",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpVal := HTTPVal{
				Method:  tt.fields.Method,
				URL:     tt.fields.URL,
				Headers: tt.fields.Headers,
			}
			got, err := httpVal.DoExecute(tt.args.requestBody, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("DoExecute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			} else {
				var responsePayload map[string]interface{}
				err = json.Unmarshal([]byte(got.(string)), &responsePayload)
				if err != nil {
					t.Errorf("DoExecute() error = %v", err)
				} else if !reflect.DeepEqual(responsePayload, tt.want) {
					t.Errorf("DoExecute() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
