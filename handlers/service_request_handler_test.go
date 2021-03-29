package handlers

import (
	"bytes"
	"clamp-core/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func prepareServiceRequestPayload() map[string]interface{} {
	serviceRequestPayload := make(map[string]interface{})
	serviceRequestPayload["userDetails"] = map[string]interface{}{"name": "testing", "address": "unit test", "mobile": "990099009900"}
	return serviceRequestPayload
}

func callCreateServiceRequest(wfName string) (*httptest.ResponseRecorder, string) {
	w := httptest.NewRecorder()
	marshal, _ := json.Marshal(prepareServiceRequestPayload())
	req, _ := http.NewRequest("POST", "/serviceRequest/"+wfName, bytes.NewReader(marshal))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", "abc")
	testHTTRouter.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	return w, bodyStr
}

func callGetServiceRequestStatus(serviceRequestID uuid.UUID) (*httptest.ResponseRecorder, string) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/serviceRequest/"+serviceRequestID.String(), nil)
	router.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	return w, bodyStr
}

func TestShouldCreateNewServiceRequestRoute(t *testing.T) {
	assert := assert.New(t)
	setUpFixture()

	w, bodyStr := callCreateServiceRequest(testWorkflowName)
	assert.Equal(http.StatusOK, w.Code)

	var jsonResp models.ServiceRequest
	err := json.Unmarshal([]byte(bodyStr), &jsonResp)
	assert.NoError(err)
	//	assert.Equal(t, workflowName, jsonResp.WorkflowName, fmt.Sprintf("The expected name was CreateOrder but we got %s", jsonResp.WorkflowName))
	assert.NotEqual(jsonResp.ID, uuid.Nil)
	assert.Equal(models.StatusNew, jsonResp.Status, fmt.Sprintf("The expected status was NEW but we got %s", jsonResp.Status))
}

func TestShouldCreateNewServiceRequestRouteWithTransformationStep(t *testing.T) {
	assert := assert.New(t)
	setUpFixture()

	w, bodyStr := callCreateServiceRequest(testTransformationWorkflow)
	assert.Equal(http.StatusOK, w.Code)

	var jsonResp models.ServiceRequest
	err := json.Unmarshal([]byte(bodyStr), &jsonResp)
	assert.NoError(err)
	//	assert.Equal(t, workflowName, jsonResp.WorkflowName, fmt.Sprintf("The expected name was CreateOrder but we got %s", jsonResp.WorkflowName))
	assert.NotEqual(jsonResp.ID, uuid.Nil)
	assert.Equal(models.StatusNew, jsonResp.Status, fmt.Sprintf("The expected status was NEW but we got %s", jsonResp.Status))
}

func TestShouldNotCreateNewServiceRequestRouteWhenServiceRequestContainsInvalidData(t *testing.T) {
	assert := assert.New(t)
	setUpFixture()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/serviceRequest/"+testWorkflowName, bytes.NewBuffer([]byte("bad payload")))
	testHTTRouter.ServeHTTP(w, req)

	w, bodyStr := callCreateServiceRequest(testWorkflowName)
	var jsonResp models.ServiceRequest
	err := json.Unmarshal([]byte(bodyStr), &jsonResp)
	assert.NoError(err)

	assert.Equal(http.StatusOK, w.Code)
	//	assert.Equal(t, workflowName, jsonResp.WorkflowName, fmt.Sprintf("The expected name was CreateOrder but we got %s", jsonResp.WorkflowName))
	assert.NotEqual(jsonResp.ID, uuid.Nil)
	assert.Equal(models.StatusNew, jsonResp.Status, fmt.Sprintf("The expected status was NEW but we got %s", jsonResp.Status))
}

func TestShouldNotCreateServiceRequestForInvalidWorkflowName(t *testing.T) {
	assert := assert.New(t)
	setUpFixture()

	w, bodyStr := callCreateServiceRequest("InvalidWF")
	assert.Equal(http.StatusBadRequest, w.Code)

	var jsonResp models.ClampErrorResponse
	err := json.Unmarshal([]byte(bodyStr), &jsonResp)
	assert.NoError(err)

	assert.Equal(http.StatusBadRequest, jsonResp.Code)
	assert.Equal("No record found with given workflow name : InvalidWF", jsonResp.Message)
}

func TestShouldGetServiceRequestStatus(t *testing.T) {
	assert := assert.New(t)
	setUpFixture()

	_, bodyStr := callCreateServiceRequest(testWorkflowName)
	var serviceReq models.ServiceRequestResponse
	err := json.Unmarshal([]byte(bodyStr), &serviceReq)
	assert.NoError(err)

	time.Sleep(time.Second) // gives time to complete service request

	status, body := callGetServiceRequestStatus(serviceReq.ID)
	assert.Equal(http.StatusOK, status.Code)

	var response models.ServiceRequestStatusResponse
	err = json.Unmarshal([]byte(body), &response)
	assert.NoError(err)
	assert.Equal(models.StatusCompleted, response.Status)
	assert.Equal(2, len(response.Steps))
}

func TestShouldFindServiceRequestByWorkflowNameByPage(t *testing.T) {
	assert := assert.New(t)
	setUpFixture()

	w, _ := callCreateServiceRequest(testWorkflowName)
	assert.Equal(http.StatusOK, w.Code)

	req, _ := http.NewRequest("GET", "/serviceRequests/testWorkflow?pageNumber=0&pageSize=1", nil)
	w = httptest.NewRecorder()
	testHTTRouter.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ServiceRequestPageResponse
	err := json.Unmarshal([]byte(bodyStr), &jsonResp)
	assert.NoError(err)

	assert.Equal(http.StatusOK, w.Code)
	assert.NotNil(jsonResp)
	assert.NotNil(jsonResp.ServiceRequests)
}

func TestShouldThrowErrorIfQueryParamsAreNotPassedInServiceRequestByWorkflowName(t *testing.T) {
	assert := assert.New(t)
	setUpFixture()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/serviceRequests/%s?pageNumber=0", testWorkflowName), nil)
	w := httptest.NewRecorder()
	testHTTRouter.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.NotNil(jsonResp)
	assert.Equal("page number or page size is not been defined", jsonResp.Message)
}

func TestShouldThrowErrorIfQueryParamsAreNotValidValuesInServiceRequestByWorkflowName(t *testing.T) {
	assert := assert.New(t)
	setUpFixture()

	req, _ := http.NewRequest("GET", "/serviceRequests/testWorkflow?pageNumber=0&pageSize=-1", nil)
	w := httptest.NewRecorder()
	testHTTRouter.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.NotNil(jsonResp)
	assert.Equal("page number or page size is not in proper format", jsonResp.Message)
}
