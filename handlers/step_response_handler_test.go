package handlers

import (
	"bytes"
	"clamp-core/models"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordStepResponse(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	res := models.ResumeStepResponse{
		ServiceRequestId: uuid.UUID{},
		StepId:           0,
		Payload:          nil,
		Errors:           models.ClampErrorResponse{},
		StepProcessed:    false,
	}
	workflowJsonReg, _ := json.Marshal(res)
	requestReader := bytes.NewReader(workflowJsonReg)

	req, _ := http.NewRequest("POST", "/stepResponse", requestReader)
	router.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampSuccessResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "success", jsonResp.Message)
}
