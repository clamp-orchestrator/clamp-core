package handlers

import (
	"bytes"
	"clamp-core/executors"
	"clamp-core/models"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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
	httpVal := executors.HttpVal{
		Method:  "GET",
		Url:     "http: //35.166.176.234:3333/api/v1/user",
		Headers: "",
	}
	steps[0] = models.Step{
		Name:      "firstStep",
		StepType:  "SYNC",
		Mode:      "HTTP",
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
	workflowReg := setUpWorkflowRequest()
	router := setupRouter()
	w := httptest.NewRecorder()
	workflowJsonReg, _ := json.Marshal(workflowReg)
	requestReader := bytes.NewReader(workflowJsonReg)

	req, _ := http.NewRequest("POST", "/workflow", requestReader)
	router.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.Workflow
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(t, 200, w.Code)
	assert.NotNil(t, jsonResp)
	assert.NotNil(t, jsonResp)
	assert.Equal(t, workflowReg.Name, jsonResp.Name, fmt.Sprintf("The expected name was %s but we got %s", workflowReg.Name, jsonResp.Name))
	assert.Equal(t, workflowDescription, jsonResp.Description, fmt.Sprintf("The expected description was Testing workflow service but the value was %s", jsonResp.Description))
	assert.NotNil(t, jsonResp.Steps)
}

func TestShouldThrowErrorIfNameFieldsIsNotPresent(t *testing.T) {
	workflowReg := setUpWorkflowRequest()
	workflowReg.Name = ""
	router := setupRouter()
	w := httptest.NewRecorder()
	workflowJsonReg, _ := json.Marshal(workflowReg)
	requestReader := bytes.NewReader(workflowJsonReg)

	req, _ := http.NewRequest("POST", "/workflow", requestReader)
	router.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)
	assert.Equal(t, http.StatusBadRequest, jsonResp.Code)
	assert.Equal(t, "Key: 'Workflow.Name' Error:Field validation for 'Name' failed on the 'required' tag", jsonResp.Message)
}

func TestShouldThrowErrorIfStepsAreNotPresent(t *testing.T) {
	workflowReg := setUpWorkflowRequest()
	workflowReg.Steps = nil
	router := setupRouter()
	w := httptest.NewRecorder()
	workflowJsonReg, _ := json.Marshal(workflowReg)
	requestReader := bytes.NewReader(workflowJsonReg)

	req, _ := http.NewRequest("POST", "/workflow", requestReader)
	router.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)
	assert.Equal(t, http.StatusBadRequest, jsonResp.Code)
	assert.Equal(t, "Key: 'Workflow.Steps' Error:Field validation for 'Steps' failed on the 'required' tag", jsonResp.Message)

	workflowReg = setUpWorkflowRequest()
	workflowReg.Steps = []models.Step{}
	router = setupRouter()
	w = httptest.NewRecorder()
	workflowJsonReg, _ = json.Marshal(workflowReg)
	requestReader = bytes.NewReader(workflowJsonReg)

	req, _ = http.NewRequest("POST", "/workflow", requestReader)
	router.ServeHTTP(w, req)

	bodyStr = w.Body.String()
	json.Unmarshal([]byte(bodyStr), &jsonResp)
	assert.Equal(t, http.StatusBadRequest, jsonResp.Code)
	assert.Equal(t, "Key: 'Workflow.Steps' Error:Field validation for 'Steps' failed on the 'gt' tag", jsonResp.Message)
}

func TestShouldThrowErrorIfStepRequiredFieldsAreNotPresent(t *testing.T) {
	workflowReg := setUpWorkflowRequest()
	workflowReg.Steps[0].Name = ""
	workflowReg.Steps[0].Mode = "HTTP"
	router := setupRouter()
	w := httptest.NewRecorder()
	workflowJsonReg, _ := json.Marshal(workflowReg)
	requestReader := bytes.NewReader(workflowJsonReg)

	req, _ := http.NewRequest("POST", "/workflow", requestReader)
	router.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)
	assert.Equal(t, http.StatusBadRequest, jsonResp.Code)
	errorMessages := strings.Split(jsonResp.Message, "\n")
	assert.Equal(t, "Key: 'Workflow.Steps[0].Name' Error:Field validation for 'Name' failed on the 'required' tag", errorMessages[0])
}
