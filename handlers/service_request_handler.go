package handlers

import (
	"clamp-core/models"
	"clamp-core/services"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

var (
	serviceRequestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "create_service_request_handler_counter",
		Help: "The total number of service requests created",
	}, []string{"workflow_name"})
	serviceRequestHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "create_service_request_handler_histogram",
		Help: "The total number of service requests created",
	})
	serviceRequestByIDCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "get_service_request_handler_by_id_counter",
		Help: "The total number of service requests enquired",
	}, []string{"service_request_id"})
	serviceRequestByIDHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "get_service_request_handler_by_id_histogram",
		Help: "The total number of service requests enquired",
	})
)

// Create Service Request godoc
// @Summary Create a service request
// @Description Create a service request and get service request id
// @Accept json
// @Produce json
// @Param workflowname path string true "Workflow Name"
// @Param serviceRequestPayload body string true "Service Request Payload"
// @Success 200 {object} models.ServiceRequestResponse
// @Failure 400 {object} models.ClampErrorResponse
// @Failure 404 {object} models.ClampErrorResponse
// @Failure 500 {object} models.ClampErrorResponse
// @Router /serviceRequest/{workflowname} [post]
func createServiceRequestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		log.Debug("Create service request handler")
		workflowName := c.Param("workflowName")
		serviceRequestCounter.WithLabelValues(workflowName).Inc()
		_, err := services.FindWorkflowByName(workflowName)

		if err != nil {
			errorResponse := models.CreateErrorResponse(http.StatusBadRequest, "No record found with given workflow name : "+workflowName)
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}

		requestPayload := readRequestPayload(c)
		// Create new service request
		serviceReq := models.NewServiceRequest(workflowName, requestPayload)
		serviceReq, err = services.SaveServiceRequest(serviceReq)
		if err != nil {
			errorResponse := models.CreateErrorResponse(http.StatusBadRequest, err.Error())
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		requestHeaders := readRequestHeadersAndSetInServiceRequest(c)
		serviceReq.RequestHeaders = requestHeaders
		services.AddServiceRequestToChannel(serviceReq)
		response := prepareServiceRequestResponse(serviceReq)
		serviceRequestHistogram.Observe(time.Since(startTime).Seconds())
		c.JSON(http.StatusOK, response)
	}
}

func readRequestHeadersAndSetInServiceRequest(c *gin.Context) string {
	var serviceRequestHeaders string
	for key, value := range c.Request.Header {
		serviceRequestHeaders += key + ":" + value[0] + ";"
	}
	// Setting Request Headers if it exists
	if serviceRequestHeaders != "" {
		log.Debugf("Service Request Headers ====> %s", serviceRequestHeaders)
		return serviceRequestHeaders
	}
	return serviceRequestHeaders
}

func prepareServiceRequestResponse(serviceReq *models.ServiceRequest) models.ServiceRequestResponse {
	response := models.ServiceRequestResponse{
		URL:    "/serviceRequest/" + serviceReq.ID.String(),
		Status: serviceReq.Status,
		ID:     serviceReq.ID,
	}
	return response
}

func readRequestPayload(c *gin.Context) map[string]interface{} {
	var payload map[string]interface{}
	if c.Request.Body != nil {
		data, _ := ioutil.ReadAll(c.Request.Body)
		_ = json.Unmarshal(data, &payload)
		log.Debug("Request Body", payload)
	}
	return payload
}

// Get Service Request By ID godoc
// @Summary Get service request details by service request id
// @Description Get service request by service request id
// @Accept json
// @Produce json
// @Param serviceRequestId path string true "Service Request ID"
// @Success 200 {object} models.ServiceRequestStatusResponse
// @Failure 400 {object} models.ClampErrorResponse
// @Failure 404 {object} models.ClampErrorResponse
// @Failure 500 {object} models.ClampErrorResponse
// @Router /serviceRequest/{serviceRequestId} [get]
func getServiceRequestStatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		serviceRequestID := c.Param("serviceRequestId")

		serviceRequest, err := services.FindServiceRequestByID(uuid.MustParse(serviceRequestID))
		if err != nil {
			c.JSON(http.StatusBadRequest, models.CreateErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}
		serviceRequestByIDCounter.WithLabelValues(serviceRequestID).Inc()
		workflow, err := services.FindWorkflowByName(serviceRequest.WorkflowName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.CreateErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}
		stepsStatues, err := services.FindStepStatusByServiceRequestID(uuid.MustParse(serviceRequestID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.CreateErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}
		stepsStatusResponse := services.PrepareStepStatusResponse(uuid.MustParse(serviceRequestID), workflow, stepsStatues)
		// TODO - handle error scenario. Currently it is always 200 ok
		serviceRequestByIDHistogram.Observe(time.Since(startTime).Seconds())
		c.JSON(http.StatusOK, stepsStatusResponse)
	}
}

func findServiceRequestByWorkflowNameHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debug("Get service request by workflow name handler")
		pageSizeStr := c.Query("pageSize")
		pageNumberStr := c.Query("pageNumber")
		if pageSizeStr == "" || pageNumberStr == "" {
			err := errors.New("page number or page size is not been defined")
			prepareErrorResponse(err, c)
			return
		}
		pageNumber, pageNumberErr := strconv.Atoi(pageNumberStr)
		pageSize, pageSizeErr := strconv.Atoi(pageSizeStr)
		if pageNumberErr != nil || pageSizeErr != nil || pageSize < 0 || pageNumber < 0 {
			err := errors.New("page number or page size is not in proper format")
			prepareErrorResponse(err, c)
			return
		}
		workflowName := c.Param("workflowName")
		serviceRequests, err := services.FindServiceRequestByWorkflowName(workflowName, pageNumber, pageSize)
		if err != nil {
			prepareErrorResponse(err, c)
			return
		}
		c.JSON(http.StatusOK, prepareServiceRequestsResponse(serviceRequests, pageNumber, pageSize))
	}
}

func prepareServiceRequestsResponse(
	serviceRequests []*models.ServiceRequest, pageNumber int, pageSize int) models.ServiceRequestPageResponse {
	response := models.ServiceRequestPageResponse{
		ServiceRequests: serviceRequests,
		PageNumber:      pageNumber,
		PageSize:        pageSize,
	}
	return response
}
