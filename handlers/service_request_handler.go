package handlers

import (
	. "clamp-core/models"
	"clamp-core/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func createServiceRequestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowName := c.Param("workflow")
		workflow, err := services.FindWorkflowByName(workflowName)
		if err != nil {
			c.JSON(http.StatusBadRequest, "No record found with given workflow name : "+workflowName)
			return
		}
		fmt.Println("Loaded workflow -", workflow)
		// Create new service request
		serviceReq := NewServiceRequest(workflowName)
		serviceReq, _ = services.SaveServiceRequest(serviceReq)
		services.AddServiceRequestToChannel(serviceReq)
		//TODO - handle error scenario. Currently it is always 200 ok
		c.JSON(http.StatusOK, serviceReq)
	}
}

func getServiceRequestStatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceRequestId := c.Param("serviceRequestId")

		var serviceRequestRes *ServiceRequest
		serviceRequestRes, _ = services.FindServiceRequestByID(uuid.MustParse(serviceRequestId))
		//TODO - handle error scenario. Currently it is always 200 ok
		c.JSON(http.StatusOK, serviceRequestRes)
	}
}
