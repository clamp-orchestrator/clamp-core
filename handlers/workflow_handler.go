package handlers

import (
	"clamp-core/models"
	"clamp-core/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)
// Create a Workflow godoc
// @Summary Create workflow for execution
// @Description Create workflow for sequential execution
// @Accept json
// @Consume json
// @Param workflowPayload body models.Workflow true "Workflow Definition Payload"
// @Success 200 {object} models.Workflow
// @Failure 400 {object} models.ClampErrorResponse
// @Failure 404 {object} models.ClampErrorResponse
// @Failure 500 {object} models.ClampErrorResponse
// @Router /workflow [post]
func createWorkflowHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create new Service Workflow
		var workflowReq models.Workflow
		err := c.ShouldBindJSON(&workflowReq)
		if err != nil {
			errorResponse := models.CreateErrorResponse(http.StatusBadRequest, err.Error())
			log.Println(err)
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		log.Printf("Create workflow request : %v \n", workflowReq)
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
// Fetch workflow details By Workflow Name godoc
// @Summary Fetch workflow details By Workflow Name
// @Description Fetch workflow details By Workflow Name
// @Accept json
// @Produce json
// @Param workflowname path string true "workflow name"
// @Success 200 {object} models.Workflow
// @Failure 400 {object} models.Workflow
// @Failure 404 {object} models.Workflow
// @Failure 500 {object} models.Workflow
// @Router /workflow/{workflowname} [get]
func fetchWorkflowBasedOnWorkflowName() gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowName := c.Param("workflow")
		result, _ := services.FindWorkflowByName(workflowName)
		//TODO - handle error scenario. Currently it is always 200 ok
		c.JSON(http.StatusOK, result)
	}
}
