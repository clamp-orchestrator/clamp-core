package executors

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpVal_DoExecute(t *testing.T) {
	assert := assert.New(t)

	httpRequestBody := map[string]interface{}{"k": "v"}
	httpResponseBody := map[string]interface{}{"id": "1234", "name": "ABC", "email": "abc@sahaj.com", "org": "sahaj"}

	var requestBody interface{}
	var requestHeader http.Header

	httpTestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "Not found")
			return
		}

		json.NewDecoder(r.Body).Decode(&requestBody)
		requestHeader = r.Header

		json.NewEncoder(w).Encode(httpResponseBody)
	}))

	defer httpTestServer.Close()

	type args struct {
		requestBody interface{}
	}
	tests := []struct {
		name              string
		fields            HTTPVal
		args              args
		wantRequestBody   interface{}
		wantRequestHeader http.Header
		wantResponseBody  interface{}
		wantErrMsg        string
	}{
		{
			name: "TestShouldExecuteHTTPStep",
			fields: HTTPVal{
				Method:  "GET",
				URL:     httpTestServer.URL,
				Headers: "Content-Type:application/json;X-Token:abc",
			},
			args: args{
				requestBody: httpRequestBody,
			},
			wantRequestBody:   httpRequestBody,
			wantRequestHeader: http.Header{"Content-Type": {"application/json"}, "X-Token": {"abc"}},
			wantResponseBody:  map[string]interface{}{"id": "1234", "name": "ABC", "email": "abc@sahaj.com", "org": "sahaj"},
			wantErrMsg:        "",
		},
		{
			name: "TestShouldThrowErrorWhileExecutingStep",
			fields: HTTPVal{
				Method:  "GET",
				URL:     httpTestServer.URL + "/non_existent_path",
				Headers: "",
			},
			wantErrMsg: "Not found",
		},
		{
			name: "TestShouldThrowErrorForHTTPStep",
			fields: HTTPVal{
				Method:  "GET",
				URL:     "http://localhost:3333/api/v1/user",
				Headers: "",
			},
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
			if tt.wantErrMsg != "" {
				assert.Contains(err.Error(), tt.wantErrMsg)
			} else {
				assert.NoError(err)

				var responseBody map[string]interface{}
				err = json.Unmarshal([]byte(got.(string)), &responseBody)
				if err != nil {
					t.Errorf("DoExecute() error = %v", err)
				}

				if tt.wantRequestBody != nil {
					assert.Equal(tt.wantRequestBody, requestBody)
				}

				if tt.wantRequestHeader != nil {
					for k, v := range tt.wantRequestHeader {
						v2, ok := requestHeader[k]
						assert.True(ok)
						assert.Equal(v, v2)
					}
				}

				if tt.wantResponseBody != nil {
					assert.Equal(tt.wantResponseBody, responseBody)
				}
			}
		})
	}
}

func TestHttpVal_PopulateRequestHeaders(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		name           string
		httpValHeaders string
		wantHTTPHeader http.Header
	}{
		{
			name:           "SingleHeader",
			httpValHeaders: "Content-Type:application/json",
			wantHTTPHeader: http.Header{
				"Content-Type": {"application/json"},
			},
		},
		{
			name:           "MultipleHeaders",
			httpValHeaders: "Content-Type:application/json;X-Header1:Value1;X-Header2:Value2",
			wantHTTPHeader: http.Header{
				"Content-Type": {"application/json"},
				"X-Header1":    {"Value1"},
				"X-Header2":    {"Value2"},
			},
		},
	}

	for i := range testCases {
		testCase := &testCases[i]

		httpHeader := make(http.Header)
		populateRequestHeaders(testCase.httpValHeaders, &httpHeader)
		assert.Equal(testCase.wantHTTPHeader, httpHeader)
	}
}
