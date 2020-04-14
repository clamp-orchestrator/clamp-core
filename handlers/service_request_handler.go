package handlers

import (
	. "clamp-core/models"
	"clamp-core/services"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
)

func createServiceRequestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Create service request handler")
		workflowName := c.Param("workflowName")
		_, err := services.FindWorkflowByName(workflowName)

		requestPayload := readRequestPayload(c)

		if err != nil {
			errorResponse := CreateErrorResponse(http.StatusBadRequest, "No record found with given workflow name : "+workflowName)
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		// Create new service request
		serviceReq := NewServiceRequest(workflowName, requestPayload)
		serviceReq, err = services.SaveServiceRequest(serviceReq)
		if err != nil {
			errorResponse := CreateErrorResponse(http.StatusBadRequest, err.Error())
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		services.AddServiceRequestToChannel(serviceReq)
		c.JSON(http.StatusOK, serviceReq)
	}
}

func readRequestPayload(c *gin.Context) map[string]interface{} {
	var payload map[string]interface{}
	if c.Request.Body != nil {
		data, _ := ioutil.ReadAll(c.Request.Body)
		json.Unmarshal(data, &payload)
		log.Println("Request Body", payload)
		return payload
	} else {
		return nil
	}
}

func getServiceRequestStatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceRequestId := c.Param("serviceRequestId")

		stepsStatues, _ := services.FindStepStatusByServiceRequestId(uuid.MustParse(serviceRequestId))
		stepsStatusResponse := services.PrepareStepStatusResponse(stepsStatues)
		//TODO - handle error scenario. Currently it is always 200 ok
		c.JSON(http.StatusOK, stepsStatusResponse)
	}
}
