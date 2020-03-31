package handlers

import (
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/serviceRequest/:workflow", createServiceRequestHandler())
	r.POST("/workflow", createWorkflowHandler())
	r.GET("/workflow/:workflow", fetchWorkflowBasedOnWorkflowName())
	return r
}

//LoadHTTPRoutes loads all HTTP api routes
func LoadHTTPRoutes() {
	r := setupRouter()
	err := r.Run()
	panic(err)
}
