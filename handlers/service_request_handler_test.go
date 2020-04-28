package handlers

import (
	"clamp-core/executors"
	"clamp-core/models"
	"clamp-core/services"
	"clamp-core/transform"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const workflowName string = "testWorkflow"

func setUp() {
	step := models.Step{
		Name:      "1",
		StepType:  "SYNC",
		Mode:      "HTTP",
		Transform: false,
		Enabled:   false,
		Val: &executors.HttpVal{
			Method:  "POST",
			Url:     "http://34.222.238.234:3333/api/v1/login",
			Headers: "",
		},
		RequestTransform: &transform.JsonTransform{
			Keys: map[string]interface{}{"a":"b"},
		},
	}

	workflow := models.Workflow{
		Name:  workflowName,
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

	w, bodyStr := callCreateServiceRequest(workflowName)
	var jsonResp models.ServiceRequest
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, workflowName, jsonResp.WorkflowName, fmt.Sprintf("The expected name was CreateOrder but we got %s", jsonResp.WorkflowName))
	assert.Equal(t, 16, len(jsonResp.ID), fmt.Sprintf("The expected length was 16 but the value was %s with length %d", jsonResp.ID, len(jsonResp.ID)))
	assert.Equal(t, models.STATUS_NEW, jsonResp.Status, fmt.Sprintf("The expected status was NEW but we got %s", jsonResp.Status))
}

func callCreateServiceRequest(wfName string) (*httptest.ResponseRecorder, string) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/serviceRequest/"+wfName, nil)
	router.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	return w, bodyStr
}

func callGetServiceRequestStatus(serviceRequestId uuid.UUID) (*httptest.ResponseRecorder, string) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/serviceRequest/"+serviceRequestId.String(), nil)
	router.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	return w, bodyStr
}

func TestShouldNotCreateServiceRequestForInvalidWorkflowName(t *testing.T) {
	setUp()
	_, bodyStr := callCreateServiceRequest("InvalidWF")
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)
	assert.Equal(t, http.StatusBadRequest, jsonResp.Code)
	assert.Equal(t, "No record found with given workflow name : InvalidWF", jsonResp.Message)
}

func TestShouldGetServiceRequestStatus(t *testing.T) {
	setUp()
	_, bodyStr := callCreateServiceRequest(workflowName)
	var serviceReq models.ServiceRequest
	json.Unmarshal([]byte(bodyStr), &serviceReq)
	time.Sleep(time.Second * 5)
	status, body := callGetServiceRequestStatus(serviceReq.ID)
	var response models.ServiceRequestStatusResponse
	json.Unmarshal([]byte(body), &response)
	assert.Equal(t, 200, status.Code)
	assert.Equal(t, models.STATUS_COMPLETED, response.Status)
	assert.Equal(t, 2, len(response.Steps))
}
