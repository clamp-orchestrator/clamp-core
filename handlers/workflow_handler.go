package handlers

import (
	"clamp-core/models"
	"clamp-core/services"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	workflowRequestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "create_workflow_request_handler_counter",
		Help: "The total number of workflow created",
	}, []string{"workflow"})
	workflowRequestHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "create_workflow_request_handler_histogram",
		Help: "The total number of service requests created",
	})
	workflowByWokflowNameCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "get_workflow_handler_by_workflow_name_counter",
		Help: "The total number of workflow fetched based on workflow name",
	}, []string{"workflow_name"})
	workflowByWokflowNameHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "get_workflow_handler_by_workflow_name_histogram",
		Help: "The total number of workflow fetched based on workflow name",
	})
)

type CustomError struct {
	StatusCode int
	Err        error
}

func (r *CustomError) Error() string {
	return fmt.Sprintf("status %d: err %v", r.StatusCode, r.Err)
}

func ErrorRequest() error {
	return &CustomError{
		StatusCode: 503,
		Err:        errors.New("unavailable"),
	}
}

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
		startTime := time.Now()
		// Create new Service Workflow
		var workflowReq models.Workflow
		workflowRequestCounter.WithLabelValues("workflow").Inc()
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
		workflowRequestHistogram.Observe(time.Since(startTime).Seconds())
		if err != nil {
			prepareErrorResponse(err, c)
			return
		}
		c.JSON(http.StatusOK, serviceFlowRes)
	}
}

func prepareErrorResponse(err error, c *gin.Context) {
	errorResponse := models.CreateErrorResponse(http.StatusBadRequest, err.Error())
	c.JSON(http.StatusBadRequest, errorResponse)
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
		startTime := time.Now()
		workflowName := c.Param("workflow")
		workflowByWokflowNameCounter.WithLabelValues(workflowName).Inc()
		result, err := services.FindWorkflowByName(workflowName)
		//TODO - handle error scenario. Currently it is always 200 ok
		workflowByWokflowNameHistogram.Observe(time.Since(startTime).Seconds())
		if err != nil {
			err := errors.New("No record exists with workflow name : " + workflowName)
			prepareErrorResponse(err, c)
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func getWorkflows() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		workflows, err := services.GetWorkflows(pageNumber, pageSize)
		if err != nil {
			prepareErrorResponse(err, c)
			return
		}
		c.JSON(http.StatusOK, prepareWorkflowResponse(workflows, pageNumber, pageSize))
	}
}

func prepareWorkflowResponse(workflows []models.Workflow, pageNumber int, pageSize int) models.WorkflowsPageResponse {
	response := models.WorkflowsPageResponse{
		Workflows:  workflows,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}
	return response
}
