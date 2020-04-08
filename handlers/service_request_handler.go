package handlers

import (
	. "clamp-core/models"
	"clamp-core/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func createServiceRequestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowName := c.Param("workflowName")
		workflow, err := services.FindWorkflowByName(workflowName)

		payload := readRequestPayload(c)

		if err != nil {
			errorResponse := CreateErrorResponse(http.StatusBadRequest, "No record found with given workflow name : "+workflowName)
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		log.Println("Loaded workflow -", workflow)
		// Create new service request
		serviceReq := NewServiceRequest(workflowName, payload)
		serviceReq, _ = services.SaveServiceRequest(serviceReq)
		services.AddServiceRequestToChannel(serviceReq)
		c.JSON(http.StatusOK, serviceReq)
	}
}

func readRequestPayload(c *gin.Context) string {
	buf := make([]byte, 1024)
	num, _ := c.Request.Body.Read(buf)
	reqBody := string(buf[0:num])
	log.Println("Request Body ", reqBody)
	return reqBody
}

func getServiceRequestStatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceRequestId := c.Param("serviceRequestId")

		var stepsStatusResponse StepsStatusResponse
		stepsStatusResponse, _ = services.FindStepStatusByServiceRequestId(uuid.MustParse(serviceRequestId))
		//TODO - handle error scenario. Currently it is always 200 ok
		c.JSON(http.StatusOK, stepsStatusResponse)
	}
}
