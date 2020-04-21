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
		var res models.AsyncStepResponse
		err := c.ShouldBindJSON(&res)
		if err != nil {
			errorResponse := models.CreateErrorResponse(http.StatusBadRequest, err.Error())
			log.Println(err)
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		log.Printf("[HTTP Consumer] : Received step completed response: %v", res)
		log.Printf("[HTTP Consumer] : Pushing step completed response to channel")
		services.AddStepResponseToResumeChannel(res)
		c.JSON(http.StatusOK, models.CreateSuccessResponse(http.StatusOK, "success"))
	}
}
