package handlers

import (
	"clamp-core/executors"
	"clamp-core/models"
	"clamp-core/repository"
	"clamp-core/services"
	"clamp-core/transform"
	"clamp-core/utils"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

const testWorkflowName string = "testWorkflow"
const testTransformationWorkflow string = "testTransformationWorkflow"

var mockDB repository.MockDB
var testHTTRouter *gin.Engine
var testHTTPServer *httptest.Server

func TestMain(m *testing.M) {
	repository.SetDB(repository.NewMemoryDB())
	gin.SetMode(gin.TestMode)

	testHTTPServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpResponseBody := map[string]interface{}{"id": "1234", "name": "ABC", "email": "abc@sahaj.com", "org": "sahaj"}
		json.NewEncoder(w).Encode(httpResponseBody)
	}))

	step := models.Step{
		Name:      "1",
		Type:      utils.StepTypeSync,
		Mode:      utils.StepModeHTTP,
		Transform: false,
		Enabled:   false,
		Val: &executors.HTTPVal{
			Method:  "POST",
			URL:     testHTTPServer.URL,
			Headers: "",
		},
	}

	workflow := models.Workflow{
		Name:  testWorkflowName,
		Steps: []models.Step{step},
	}

	_, err := services.SaveWorkflow(&workflow)
	if err != nil {
		panic(err)
	}

	step = models.Step{
		Name:      "1",
		Type:      utils.StepTypeSync,
		Mode:      utils.StepModeHTTP,
		Transform: true,
		Enabled:   false,
		RequestTransform: &transform.JSONTransform{
			Spec: map[string]interface{}{"name": "test"},
		},
		Val: &executors.HTTPVal{
			Method:  "POST",
			URL:     testHTTPServer.URL,
			Headers: "",
		},
	}

	workflow = models.Workflow{
		Name:  testTransformationWorkflow,
		Steps: []models.Step{step},
	}

	_, err = services.SaveWorkflow(&workflow)
	if err != nil {
		panic(err)
	}

	testHTTRouter = setupRouter()

	//log.SetLevel(log.DebugLevel)
	os.Exit(m.Run())
}
