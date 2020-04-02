package handlers

import (
	"bytes"
	"clamp-core/models"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"net/http/httptest"
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
	steps := make([]models.Step,1)

	steps[0] = models.Step{
		Id:      "firstStep",
		Name:    "firstStep",
		Enabled: true,
		Mode:"http",
		URL:"www.google.com",
		Transform:false,
	}
	workflow := models.Workflow{
		Name: RandStringRunes(10),
		Description: workflowDescription,
		Steps: steps,
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
