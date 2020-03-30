package handlers

import (
	. "clamp-core/models"
	"clamp-core/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func createServiceRequestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowName := c.Param("workflow")
		serviceReq := NewServiceRequest(workflowName)
		services.SaveServiceRequest(serviceReq)
		//TODO - handle error scenario. Currently it is always 200 ok
		c.JSON(http.StatusOK, serviceReq)
	}
}
