package handlers

import (
	"clamp-core/models"
	"clamp-core/services"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	resumeAsyncServiceRequestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "resume_async_service_request_handler_counter",
		Help: "The total number of async service requests resumed",
	}, []string{"resume"})
	resumeAsyncServiceRequestHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "resume_async_service_request_handler_histogram",
		Help: "The total number of async service requests resumed",
	})
)

// Http Resume Service Request API for Async Step godoc
// @Summary Http Resume Service Request API for Async Step
// @Description Http Resume Service Request API for Async Step
// @Accept json
// @Produce json
// @Param ResumeServiceRequestPayload body models.AsyncStepResponse true "Resume Service Request Payload"
// @Success 200 {object} models.ClampSuccessResponse
// @Failure 400 {object} models.ClampErrorResponse
// @Failure 404 {object} models.ClampErrorResponse
// @Failure 500 {object} models.ClampErrorResponse
// @Router /stepResponse [post]
func createStepResponseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		var res models.AsyncStepResponse
		resumeAsyncServiceRequestCounter.WithLabelValues("resume").Inc()
		err := c.ShouldBindJSON(&res)
		if err != nil {
			errorResponse := models.CreateErrorResponse(http.StatusBadRequest, err.Error())
			log.Println(err)
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		resumeServiceRequestHeaders := readRequestHeadersAndSetInServiceRequest(c)
		res.RequestHeaders = resumeServiceRequestHeaders
		log.Printf("[HTTP Consumer] : Received step completed response: %v", res)
		log.Printf("[HTTP Consumer] : Pushing step completed response to channel")
		services.AddStepResponseToResumeChannel(res)
		resumeAsyncServiceRequestHistogram.Observe(time.Since(startTime).Seconds())
		c.JSON(http.StatusOK, models.CreateSuccessResponse(http.StatusOK, "success"))
	}
}
