package handlers

import (
	"clamp-core/models"
	"clamp-core/services"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const workflowName string = "testWF"

func setUp() {
	const workflowName = "testWorkflow"
	workflow := models.Workflow{
		Name: workflowName,
	}
	resp, err := services.FindWorkflowByName(workflowName)
	fmt.Println(resp)
	if err != nil {
		services.SaveWorkflow(workflow)
	}
}

func TestCreateNewServiceRequestRoute(t *testing.T) {
	setUp()

	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/serviceRequest/"+workflowName, nil)
	router.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ServiceRequest
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, workflowName, jsonResp.WorkflowName, fmt.Sprintf("The expected name was CreateOrder but we got %s", jsonResp.WorkflowName))
	assert.Equal(t, 16, len(jsonResp.ID), fmt.Sprintf("The expected length was 16 but the value was %s with length %d", jsonResp.ID, len(jsonResp.ID)))
	assert.Equal(t, models.STATUS_NEW, jsonResp.Status, fmt.Sprintf("The expected status was NEW but we got %s", jsonResp.Status))
}
