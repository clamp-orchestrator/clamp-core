package handlers

import (
	"bytes"
	"clamp-core/executors"
	"clamp-core/models"
	"clamp-core/services"
	"clamp-core/transform"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const workflowName string = "testWorkflow"
const transformationWorkflowName string = "transformWorkflow"

func TestShouldCreateNewServiceRequestRoute(t *testing.T) {
	CreateWorkflowIfItsAlreadyDoesNotExists()

	w, bodyStr := callCreateServiceRequest(workflowName)
	var jsonResp models.ServiceRequest
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(t, 200, w.Code)
	//	assert.Equal(t, workflowName, jsonResp.WorkflowName, fmt.Sprintf("The expected name was CreateOrder but we got %s", jsonResp.WorkflowName))
	assert.Equal(t, 16, len(jsonResp.ID), fmt.Sprintf("The expected length was 16 but the value was %s with length %d", jsonResp.ID, len(jsonResp.ID)))
	assert.Equal(t, models.STATUS_NEW, jsonResp.Status, fmt.Sprintf("The expected status was NEW but we got %s", jsonResp.Status))
}

func TestShouldNotCreateNewServiceRequestRouteWithTransformationStep(t *testing.T) {
	createWorkflowWithTransformationEnabledInOneStep()

	w, bodyStr := callCreateServiceRequest(transformationWorkflowName)
	var jsonResp models.ServiceRequest
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(t, 200, w.Code)
	//	assert.Equal(t, workflowName, jsonResp.WorkflowName, fmt.Sprintf("The expected name was CreateOrder but we got %s", jsonResp.WorkflowName))
	assert.Equal(t, 16, len(jsonResp.ID), fmt.Sprintf("The expected length was 16 but the value was %s with length %d", jsonResp.ID, len(jsonResp.ID)))
	assert.Equal(t, models.STATUS_NEW, jsonResp.Status, fmt.Sprintf("The expected status was NEW but we got %s", jsonResp.Status))
}

//TODO
func TestShouldNotCreateNewServiceRequestRouteWhenServiceRequestContainsInvalidData(t *testing.T) {
	CreateWorkflowIfItsAlreadyDoesNotExists()

	w, bodyStr := callCreateServiceRequest(workflowName)
	var jsonResp models.ServiceRequest
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(t, 200, w.Code)
	//	assert.Equal(t, workflowName, jsonResp.WorkflowName, fmt.Sprintf("The expected name was CreateOrder but we got %s", jsonResp.WorkflowName))
	assert.Equal(t, 16, len(jsonResp.ID), fmt.Sprintf("The expected length was 16 but the value was %s with length %d", jsonResp.ID, len(jsonResp.ID)))
	assert.Equal(t, models.STATUS_NEW, jsonResp.Status, fmt.Sprintf("The expected status was NEW but we got %s", jsonResp.Status))
}

func TestShouldNotCreateServiceRequestForInvalidWorkflowName(t *testing.T) {
	CreateWorkflowIfItsAlreadyDoesNotExists()
	_, bodyStr := callCreateServiceRequest("InvalidWF")
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)
	assert.Equal(t, http.StatusBadRequest, jsonResp.Code)
	assert.Equal(t, "No record found with given workflow name : InvalidWF", jsonResp.Message)
}

func TestShouldGetServiceRequestStatus(t *testing.T) {
	CreateWorkflowIfItsAlreadyDoesNotExists()
	_, bodyStr := callCreateServiceRequest(workflowName)
	var serviceReq models.ServiceRequestResponse
	json.Unmarshal([]byte(bodyStr), &serviceReq)
	time.Sleep(time.Second * 5)
	status, body := callGetServiceRequestStatus(serviceReq.ID)
	var response models.ServiceRequestStatusResponse
	json.Unmarshal([]byte(body), &response)
	assert.Equal(t, 200, status.Code)
	assert.Equal(t, models.STATUS_COMPLETED, response.Status)
	assert.Equal(t, 2, len(response.Steps))
}

func TestShouldGetWorkflowNotExistsWhenGetServiceRequestStatusCalled(t *testing.T) {
	CreateWorkflowIfItsAlreadyDoesNotExists()
	_, bodyStr := callCreateServiceRequest(workflowName)
	var serviceReq models.ServiceRequestResponse
	json.Unmarshal([]byte(bodyStr), &serviceReq)
	time.Sleep(time.Second * 5)
	DeleteWorkflowIfExists()
	status, body := callGetServiceRequestStatus(serviceReq.ID)
	var response models.ClampErrorResponse
	json.Unmarshal([]byte(body), &response)
	assert.Equal(t, 500, status.Code)
	assert.Equal(t, "pg: no rows in result set", response.Message)
}

func callCreateServiceRequest(wfName string) (*httptest.ResponseRecorder, string) {
	router := setupRouter()
	w := httptest.NewRecorder()
	marshal, _ := json.Marshal(prepareServiceRequestPayload())
	req, _ := http.NewRequest("POST", "/serviceRequest/"+wfName, bytes.NewReader(marshal))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", "abc")
	router.ServeHTTP(w, req)

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

func CreateWorkflowIfItsAlreadyDoesNotExists() {
	step := models.Step{
		Name:      "1",
		Type:      "SYNC",
		Mode:      "HTTP",
		Transform: false,
		Enabled:   false,
		Val: &executors.HTTPVal{
			Method:  "POST",
			URL:     "https://run.mocky.io/v3/0590fbf8-0f1c-401c-b9df-65e98ef0385d",
			Headers: "",
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

func DeleteWorkflowIfExists() {
	err := services.DeleteWorkflowByName(workflowName)
	log.Println(err)
}

func prepareServiceRequestPayload() map[string]interface{} {
	serviceRequestPayload := make(map[string]interface{})
	serviceRequestPayload["userDetails"] = map[string]interface{}{"name": "testing", "address": "unit test", "mobile": "990099009900"}
	return serviceRequestPayload
}

func createWorkflowWithTransformationEnabledInOneStep() {
	step := models.Step{
		Name:      "1",
		Type:      "SYNC",
		Mode:      "HTTP",
		Transform: true,
		Enabled:   false,
		RequestTransform: &transform.JSONTransform{
			Spec: map[string]interface{}{"name": "test"},
		},
		Val: &executors.HTTPVal{
			Method:  "POST",
			URL:     "https://reqres.in/api/users",
			Headers: "",
		},
	}

	workflow := models.Workflow{
		Name:  transformationWorkflowName,
		Steps: []models.Step{step},
	}
	resp, err := services.FindWorkflowByName(transformationWorkflowName)
	log.Println(resp)
	if err != nil {
		services.SaveWorkflow(workflow)
	}
}

func TestShouldFindServiceRequestByWorkflowNameByPage(t *testing.T) {
	CreateWorkflowIfItsAlreadyDoesNotExists()
	router := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/serviceRequests/testWorkflow?pageNumber=1&pageSize=1", nil)
	router.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ServiceRequestPageResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(t, 200, w.Code)
	assert.NotNil(t, jsonResp)
}

func TestShouldThrowErrorIfQueryParamsAreNotPassedInServiceRequestByWorkflowName(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/serviceRequests/testWorkflow?pageNumber=1", nil)
	router.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(t, 400, w.Code)
	assert.NotNil(t, jsonResp)
	assert.Equal(t, "page number or page size is not been defined", jsonResp.Message)
}

func TestShouldThrowErrorIfQueryParamsAreNotValidValuesInServiceRequestByWorkflowName(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/serviceRequests/testWorkflow?pageNumber=0&pageSize=-1", nil)
	router.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(t, 400, w.Code)
	assert.NotNil(t, jsonResp)
	assert.Equal(t, "page number or page size is not in proper format", jsonResp.Message)
}

func TestShouldGetServiceRequestsByWorkflowNameAndBySort(t *testing.T) {
	CreateWorkflowIfItsAlreadyDoesNotExists()
	router := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/serviceRequests/testWorkflow?pageNumber=1&pageSize=1&sortBy=id:desc", nil)
	router.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ServiceRequestPageResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(t, 200, w.Code)
	assert.NotNil(t, jsonResp)
}

func TestShouldFailGetServiceRequestsByWorkflowNameWithInvalidSortFields(t *testing.T) {
	CreateWorkflowIfItsAlreadyDoesNotExists()
	router := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/serviceRequests/testWorkflow?pageNumber=1&pageSize=1&sortBy=ids:desc", nil)
	router.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ServiceRequestPageResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(t, 400, w.Code)
	assert.Nil(t, jsonResp.ServiceRequests)
}

func TestShouldFailGetServiceRequestsByWorkflowNameWithInvalidSortQuery(t *testing.T) {
	CreateWorkflowIfItsAlreadyDoesNotExists()
	router := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/serviceRequests/testWorkflow?pageNumber=1&pageSize=1&sortBy=id:", nil)
	router.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ServiceRequestPageResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(t, 400, w.Code)
	assert.Nil(t, jsonResp.ServiceRequests)
}
