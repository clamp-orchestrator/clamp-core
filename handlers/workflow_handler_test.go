package handlers

import (
	"bytes"
	"clamp-core/executors"
	"clamp-core/models"
	"clamp-core/utils"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const workflowDescription string = "Testing workflow service"

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func setUpWorkflowRequest() models.Workflow {
	steps := make([]models.Step, 1)
	httpVal := executors.HTTPVal{
		Method:  "GET",
		URL:     testHTTPServer.URL,
		Headers: "",
	}
	steps[0] = models.Step{
		Name:      "firstStep",
		Type:      utils.StepTypeSync,
		Mode:      utils.StepModeHTTP,
		Val:       httpVal,
		Transform: false,
		Enabled:   true,
	}
	workflow := models.Workflow{
		Name:        RandStringRunes(10),
		Description: workflowDescription,
		Steps:       steps,
	}
	return workflow
}

func TestCreateNewWorkflowRequestRoute(t *testing.T) {
	assert := assert.New(t)

	workflowReg := setUpWorkflowRequest()
	w := httptest.NewRecorder()
	workflowJSONReg, _ := json.Marshal(workflowReg)
	requestReader := bytes.NewReader(workflowJSONReg)

	req, _ := http.NewRequest("POST", "/workflow", requestReader)
	testHTTRouter.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.Workflow
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(http.StatusOK, w.Code)
	assert.NotNil(jsonResp)
	assert.NotNil(jsonResp)
	assert.Equal(workflowReg.Name, jsonResp.Name, fmt.Sprintf("The expected name was %s but we got %s", workflowReg.Name, jsonResp.Name))
	assert.Equal(workflowDescription, jsonResp.Description, fmt.Sprintf("The expected description was Testing workflow service but the value was %s", jsonResp.Description))
	assert.NotNil(jsonResp.Steps)
}

func TestShouldThrowErrorIfNameFieldsIsNotPresent(t *testing.T) {
	assert := assert.New(t)

	workflowReg := setUpWorkflowRequest()
	workflowReg.Name = ""

	w := httptest.NewRecorder()
	workflowJSONReg, _ := json.Marshal(workflowReg)
	requestReader := bytes.NewReader(workflowJSONReg)

	req, _ := http.NewRequest("POST", "/workflow", requestReader)
	testHTTRouter.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)
	assert.Equal(http.StatusBadRequest, jsonResp.Code)
	assert.Equal("Key: 'Workflow.Name' Error:Field validation for 'Name' failed on the 'required' tag", jsonResp.Message)
}

func TestShouldThrowErrorIfStepsAreNotPresent(t *testing.T) {
	assert := assert.New(t)

	workflowReg := setUpWorkflowRequest()
	workflowReg.Steps = nil
	w := httptest.NewRecorder()
	workflowJSONReg, _ := json.Marshal(workflowReg)
	requestReader := bytes.NewReader(workflowJSONReg)

	req, _ := http.NewRequest("POST", "/workflow", requestReader)
	testHTTRouter.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)
	assert.Equal(http.StatusBadRequest, jsonResp.Code)
	assert.Equal("Key: 'Workflow.Steps' Error:Field validation for 'Steps' failed on the 'required' tag", jsonResp.Message)

	workflowReg = setUpWorkflowRequest()
	workflowReg.Steps = []models.Step{}

	w = httptest.NewRecorder()
	workflowJSONReg, _ = json.Marshal(workflowReg)
	requestReader = bytes.NewReader(workflowJSONReg)

	req, _ = http.NewRequest("POST", "/workflow", requestReader)
	testHTTRouter.ServeHTTP(w, req)

	bodyStr = w.Body.String()
	json.Unmarshal([]byte(bodyStr), &jsonResp)
	assert.Equal(http.StatusBadRequest, jsonResp.Code)
	assert.Equal("Key: 'Workflow.Steps' Error:Field validation for 'Steps' failed on the 'gt' tag", jsonResp.Message)
}

func TestShouldThrowErrorIfStepRequiredFieldsAreNotPresent(t *testing.T) {
	assert := assert.New(t)

	workflowReg := setUpWorkflowRequest()
	workflowReg.Steps[0].Name = ""
	workflowReg.Steps[0].Mode = utils.StepModeHTTP

	w := httptest.NewRecorder()
	workflowJSONReg, _ := json.Marshal(workflowReg)
	requestReader := bytes.NewReader(workflowJSONReg)

	req, _ := http.NewRequest("POST", "/workflow", requestReader)
	testHTTRouter.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)
	assert.Equal(http.StatusBadRequest, jsonResp.Code)
	errorMessages := strings.Split(jsonResp.Message, "\n")
	assert.Equal("Key: 'Workflow.Steps[0].Name' Error:Field validation for 'Name' failed on the 'required' tag", errorMessages[0])
}

func TestShouldReturnCreatedWorkflowSuccessfullyByWorkflowNameRoute(t *testing.T) {
	assert := assert.New(t)

	workflowReg := setUpWorkflowRequest()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/workflow/"+testWorkflowName, nil)
	testHTTRouter.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.Workflow
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(http.StatusOK, w.Code)
	assert.NotNil(jsonResp)
	assert.Equal(testWorkflowName, jsonResp.Name, fmt.Sprintf("The expected name was %s but we got %s", workflowReg.Name, jsonResp.Name))
	assert.NotNil(jsonResp.Steps)
}

func TestShouldFailToReturnWorkflowIfInvalidWorkflowNameIsProvidedInTheRoute(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/workflow/"+"dummy", nil)
	testHTTRouter.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.NotNil(jsonResp)
	assert.Equal(http.StatusBadRequest, jsonResp.Code)
	assert.Equal("No record exists with workflow name : "+"dummy", jsonResp.Message)
}

func TestCreateNewWorkflowRequestShouldFailIfWorkflowNameAlreadyExistsRoute(t *testing.T) {
	assert := assert.New(t)

	workflowReg := setUpWorkflowRequest()
	workflowReg.Name = testWorkflowName
	w := httptest.NewRecorder()
	workflowJSONReg, _ := json.Marshal(workflowReg)
	requestReader := bytes.NewReader(workflowJSONReg)

	req, _ := http.NewRequest("POST", "/workflow", requestReader)
	testHTTRouter.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var errorJSONResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &errorJSONResp)
	assert.Equal(http.StatusBadRequest, w.Code)
	assert.NotNil(errorJSONResp.Code)
	assert.NotNil(errorJSONResp.Message)
}

func TestShouldGetAllWorkflowsByPage(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/workflows?pageNumber=1&pageSize=1", nil)
	testHTTRouter.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.WorkflowsPageResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(http.StatusOK, w.Code)
	assert.NotNil(jsonResp)
	assert.Equal(1, len(jsonResp.Workflows), fmt.Sprintf("The expected number of records was %d but we got %d", 1, len(jsonResp.Workflows)))
	assert.True(jsonResp.TotalWorkflowsCount > 0, "The total workflow count is less than 0")
}

func TestShouldThrowErrorIfQueryParamsAreNotPassedInGetAllWorkflows(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/workflows?pageNumber=1", nil)
	testHTTRouter.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.NotNil(jsonResp)
	assert.Equal("page number or page size has not been defined", jsonResp.Message)
}

func TestShouldThrowErrorIfPageNumberIsLessThanOne(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/workflows?pageNumber=0&pageSize=1", nil)
	testHTTRouter.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.NotNil(jsonResp)
	assert.Equal("page number or page size is not in proper format", jsonResp.Message)
}

func TestShouldThrowErrorIfPageSizeIsLessThanOne(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/workflows?pageNumber=1&pageSize=0", nil)
	testHTTRouter.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.NotNil(jsonResp)
	assert.Equal("page number or page size is not in proper format", jsonResp.Message)
}

func TestShouldThrowErrorIfQueryParamsAreNotValidValuesInGetAllWorkflows(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/workflows?pageNumber=1&pageSize=-1", nil)
	testHTTRouter.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.NotNil(jsonResp)
	assert.Equal("page number or page size is not in proper format", jsonResp.Message)
}

func TestShouldThrowErrorIfSortByStringIsNotInTheRightFormat(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	sortByString := `id,`
	urlEncodedSortValue := url.QueryEscape(sortByString)
	req, _ := http.NewRequest("GET", "/workflows?pageNumber=1&pageSize=1&sortBy="+urlEncodedSortValue, nil)
	testHTTRouter.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.NotNil(jsonResp)
	assert.Equal("unsupported value provided for sortBy query", jsonResp.Message)
}

func TestShouldThrowErrorIfSortContainsInvalidFields(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	sortByString := `updatedate:ASC`
	urlEncodedSortValue := url.QueryEscape(sortByString)
	req, _ := http.NewRequest("GET", "/workflows?pageNumber=1&pageSize=1&sortBy="+urlEncodedSortValue, nil)
	testHTTRouter.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.NotNil(jsonResp)
	assert.Equal("unsupported value provided for sortBy query", jsonResp.Message)
}
