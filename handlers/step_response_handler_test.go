package handlers

import (
	"bytes"
	"clamp-core/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRecordStepResponse(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	res := models.AsyncStepResponse{
		ServiceRequestID: uuid.UUID{},
		StepID:           0,
		Response:         nil,
	}
	workflowJSONReg, _ := json.Marshal(res)
	requestReader := bytes.NewReader(workflowJSONReg)

	req, _ := http.NewRequest("POST", "/stepResponse", requestReader)
	router.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampSuccessResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "success", jsonResp.Message)
}

func TestShouldReturnBadRequestWhenRequestContainsInvalidDataForRecordStepResponse(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/stepResponse", nil)
	router.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "invalid request", jsonResp.Message)
}
