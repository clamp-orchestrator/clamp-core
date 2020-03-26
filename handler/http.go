package handler

import (
	"clamp-core/domain"
	"clamp-core/repository"
	"encoding/json"

	//"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createServiceRequestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowName := c.Param("workflow")
		serviceReq := domain.Create(workflowName)
		repository.SaveServiceRequest(serviceReq)
		//TODO - handle error scenario. Currently it is always 200 ok
		c.JSON(http.StatusOK, serviceReq)
	}
}

func createWorkflowRequestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody, _ := c.GetRawData()
		// Create new Service Request
		request := domain.Request{}
		json.Unmarshal([]byte(requestBody), &request)
		fmt.Printf("Operation: %v \n", request.ServiceFlow)

		serviceFlowRes := domain.CreateWorkflow(request.ServiceFlow)
		//TODO - handle error scenario. Currently it is always 200 ok
		c.JSON(http.StatusOK, serviceFlowRes)
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/serviceRequest/:workflow", createServiceRequestHandler())
	r.POST("/workflow", createWorkflowRequestHandler())
	return r
}

//LoadHTTPRoutes loads all HTTP api routes
func LoadHTTPRoutes() {
	r := setupRouter()
	r.Run()
}
