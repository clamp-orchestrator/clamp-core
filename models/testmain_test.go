package models

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var testHTTPServer *httptest.Server

func TestMain(m *testing.M) {
	testHTTPServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpResponseBody := map[string]interface{}{"id": "1234", "name": "ABC", "email": "abc@sahaj.com", "org": "sahaj"}
		json.NewEncoder(w).Encode(httpResponseBody)
	}))

	os.Exit(m.Run())
}
