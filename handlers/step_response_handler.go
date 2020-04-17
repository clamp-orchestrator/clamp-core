package handlers

import (
	"clamp-core/models"
	"clamp-core/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func createStepResponseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.ResumeStepResponse
		err := c.ShouldBindJSON(&request)
		if err != nil {
			errorResponse := models.CreateErrorResponse(http.StatusBadRequest, err.Error())
			log.Println(err)
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		log.Printf("[HTTP Consumer] : Received step completed response: %v", request)
		log.Printf("[HTTP Consumer] : Pushing step completed response to channel")
		services.AddAsyncResumeStepExecutionRequestToChannel(request)
		c.JSON(http.StatusOK, models.CreateSuccessResponse(http.StatusOK, "success"))
	}
}
