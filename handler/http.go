package handler

import (
	"clamp-core/servicerequest"
	"github.com/gin-gonic/gin"
	"net/http"
)

func createServiceRequestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowName := c.Param("workflow")
		servicerequest.Create(workflowName)

		c.JSON(http.StatusOK, gin.H{
			"message": "ok!",
		})
	}
}

//LoadHTTPRoutes loads all HTTP api routes
func LoadHTTPRoutes() {
	r := gin.Default()

	r.POST("/serviceRequest/:workflow", createServiceRequestHandler())

	r.Run()
}
