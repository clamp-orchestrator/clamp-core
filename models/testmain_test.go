package models

import (
	"clamp-core/repository"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var mockDB repository.MockDB
var testHTTPServer *httptest.Server

func TestMain(m *testing.M) {
	testHTTPServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpResponseBody := map[string]interface{}{"id": "1234", "name": "ABC", "email": "abc@sahaj.com", "org": "sahaj"}
		json.NewEncoder(w).Encode(httpResponseBody)
	}))

	repository.SetDB(&mockDB)

	os.Exit(m.Run())
}
