package handler

import (
	"clamp-core/servicerequest"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createServiceRequestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowName := c.Param("workflow")
		serviceReq := servicerequest.Create(workflowName)

		//TODO - handle error scenario. Currently it is always 200 ok
		c.JSON(http.StatusOK, serviceReq)
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/serviceRequest/:workflow", createServiceRequestHandler())
	return r
}

//LoadHTTPRoutes loads all HTTP api routes
func LoadHTTPRoutes() {
	r := setupRouter()
	r.Run()
}
