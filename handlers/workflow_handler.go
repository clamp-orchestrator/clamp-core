package handlers

import (
	"clamp-core/models"
	"clamp-core/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func createWorkflowHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create new Service Workflow
		var workflowReq models.Workflow
		err := c.ShouldBindJSON(&workflowReq)
		if err != nil {
			errorResponse := models.CreateErrorResponse(http.StatusBadRequest, err.Error())
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		log.Printf("Create workflowReq workflowReq : %v \n", workflowReq)
		serviceFlowRes := models.CreateWorkflow(workflowReq)
		serviceFlowRes, err = services.SaveWorkflow(serviceFlowRes)
		if err != nil {
			errorResponse := models.CreateErrorResponse(http.StatusBadRequest, err.Error())
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		c.JSON(http.StatusOK, serviceFlowRes)
	}
}

func fetchWorkflowBasedOnWorkflowName() gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowName := c.Param("workflow")
		result, _ := services.FindWorkflowByName(workflowName)
		//TODO - handle error scenario. Currently it is always 200 ok
		c.JSON(http.StatusOK, result)
	}
}
