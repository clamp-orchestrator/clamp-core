package handlers

import (
	"clamp-core/executors"
	"clamp-core/models"
	"clamp-core/services"
	"encoding/json"
	"fmt"
	"log"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const workflowName string = "testWF"

func setUp() {
	step := models.Step{
		Id:        "1",
		Name:      "1",
		StepType:  "SYNC",
		Mode:      "HTTP",
		Transform: false,
		Enabled:   false,
		Val: &executors.HttpVal{
			Method:  "POST",
			Url:     "http://35.166.176.234:3333/api/v1/login",
			Headers: "",
		},
	}

	workflow := models.Workflow{
		Name: workflowName,
		Steps: []models.Step{step},
	}
	resp, err := services.FindWorkflowByName(workflowName)
	log.Println(resp)
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

func TestShouldNotCreateServiceRequestForInvalidWorkflowName(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/serviceRequest/InvalidWF", nil)
	router.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)
	assert.Equal(t, http.StatusBadRequest, jsonResp.Code)
	assert.Equal(t, "No record found with given workflow name : InvalidWF", jsonResp.Message)
}
