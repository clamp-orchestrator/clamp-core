package handlers

import (
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/metrics", prometheusHandler())
	r.POST("/serviceRequest/:workflowName", createServiceRequestHandler())
	r.GET("/serviceRequest/:serviceRequestId", getServiceRequestStatusHandler())
	r.POST("/workflow", createWorkflowHandler())
	r.GET("/workflow/:workflow", fetchWorkflowBasedOnWorkflowName())
	r.POST("/stepResponse", createStepResponseHandler())
	return r
}

//LoadHTTPRoutes loads all HTTP api routes
func LoadHTTPRoutes() {
	r := setupRouter()
	err := r.Run()
	panic(err)
}
