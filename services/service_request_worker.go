package services

import (
	"clamp-core/models"
	"clamp-core/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

var (
	completedServiceRequestCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "completed_service_request_handler_counter",
		Help: "The total number of service requests completed",
	})
	failedServiceRequestCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "failed_service_request_handler_counter",
		Help: "The total number of service requests failed",
	})
)
var (
	serviceRequestChannel chan models.ServiceRequest
	singletonOnce         sync.Once
)

func createServiceRequestChannel() chan models.ServiceRequest {
	singletonOnce.Do(func() {
		serviceRequestChannel = make(chan models.ServiceRequest, utils.ServiceRequestChannelSize)
	})
	return serviceRequestChannel
}

func init() {
	createServiceRequestChannel()
	createServiceRequestWorkers()
}

func createServiceRequestWorkers() {
	for i := 0; i < utils.ServiceRequestWorkersSize; i++ {
		go worker(i, serviceRequestChannel)
	}
}

func worker(workerID int, serviceReqChan <-chan models.ServiceRequest) {
	prefix := fmt.Sprintf("[WORKER_%d] ", workerID)
	prefix = fmt.Sprintf("%15s", prefix)
	log.Infof("%s : Started listening to service request channel", prefix)
	for serviceReq := range serviceReqChan {
		serviceReqCopy := serviceReq
		executeWorkflow(&serviceReqCopy, prefix)
	}
}

func executeWorkflow(serviceReq *models.ServiceRequest, prefix string) {
	prefix = fmt.Sprintf("%s [REQUEST_ID: %s]", prefix, serviceReq.ID)
	log.Debugf("%s Started processing service request id %s", prefix, serviceReq.ID)
	start := time.Now()
	workflow, err := FindWorkflowByName(serviceReq.WorkflowName)
	if err == nil {
		lastStep := workflow.Steps[len(workflow.Steps)-1]
		if serviceReq.CurrentStepID == 0 || serviceReq.CurrentStepID != lastStep.ID {
			status := executeWorkflowSteps(workflow, prefix, serviceReq)
			if status == models.StatusCompleted {
				completedServiceRequestCounter.Inc()
				elapsed := time.Since(start)
				log.Debugf("%s Completed processing service request id %s in %s", prefix, serviceReq.ID, elapsed)
			} else if status == models.StatusFailed {
				failedServiceRequestCounter.Inc()
			}
		} else {
			log.Debugf("%s All steps are executed for service request id: %s", prefix, serviceReq.ID)
		}
	}
}

func catchErrors(prefix string, requestID uuid.UUID) {
	if r := recover(); r != nil {
		log.Error("[ERROR]", r)
		log.Errorf("%s Failed processing service request id %s", prefix, requestID)
	}
}

func executeWorkflowSteps(workflow *models.Workflow, prefix string, serviceRequest *models.ServiceRequest) models.Status {
	stepRequestPayload := serviceRequest.Payload
	lastStepExecuted := serviceRequest.CurrentStepID
	executeStepsFromIndex := 0
	if lastStepExecuted > 0 {
		executeStepsFromIndex = lastStepExecuted
		log.Debugf("%s Skipping steps till step id %d", prefix, executeStepsFromIndex)
	}
	requestContext := CreateRequestContext(workflow, serviceRequest)
	// prepare request context for async steps

	if executeStepsFromIndex > 0 {
		EnhanceRequestContextWithExecutedSteps(&requestContext)
	}

	for i := range workflow.Steps[executeStepsFromIndex:] {
		step := &workflow.Steps[i]
		ComputeRequestToCurrentStepInContext(workflow, step, &requestContext, executeStepsFromIndex+i, stepRequestPayload)
		err := ExecuteWorkflowStep(step, requestContext, prefix)
		if !err.IsNil() {
			return models.StatusFailed
		}
		if !requestContext.StepsContext[step.Name].StepSkipped && step.Type == utils.StepTypeAsync {
			log.Debugf("%s : Pushed to sleep mode until response for step - %s is received", prefix, step.Name)
			return models.StatusPaused
		}
	}
	return models.StatusCompleted
}

// TODO: replace prefix with other standard way like MDC
func ExecuteWorkflowStep(step *models.Step, requestContext models.RequestContext, prefix string) models.ClampErrorResponse {
	serviceRequestID := requestContext.ServiceRequestID
	workflowName := requestContext.WorkflowName
	stepRequest := requestContext.StepsContext[step.Name].Request

	defer catchErrors(prefix, serviceRequestID)

	requestContext.SetStepRequestToContext(step.Name, stepRequest)

	stepStartTime := time.Now()
	stepStatus := &models.StepsStatus{
		ServiceRequestID: serviceRequestID,
		WorkflowName:     workflowName,
		StepName:         step.Name,
		Payload: models.Payload{
			Request:  stepRequest,
			Response: nil,
		},
		StepID: step.ID,
	}

	// TODO Condition should be checked on transformed request or original request? Based on that this section needs to be altered
	if step.Transform {
		transform, transformErrors := step.DoTransform(requestContext, prefix)
		if transformErrors != nil {
			log.Error("Error while transforming request payload")
			panic(transformErrors)
		}
		requestContext.SetStepRequestToContext(step.Name, transform)
		stepStatus.Payload.Request = transform
	}
	recordStepStartedStatus(stepStatus, stepStartTime)

	resp, err := step.DoExecute(requestContext, prefix)
	if err != nil {
		if step.OnFailure != nil {
			for i := range step.OnFailure {
				stepOnFailure := &step.OnFailure[i]
				ExecuteWorkflowStep(stepOnFailure, requestContext, prefix)
			}
		}
		clampErrorResponse := models.CreateErrorResponse(http.StatusBadRequest, err.Error())
		recordStepFailedStatus(stepStatus, *clampErrorResponse, stepStartTime)
		return *clampErrorResponse
	} else if step.DidStepExecute() && resp != nil && step.Type == utils.StepTypeSync {
		log.Debugf("%s Step response received: %s", prefix, resp.(string))
		var responsePayload map[string]interface{}
		_ = json.Unmarshal([]byte(resp.(string)), &responsePayload)
		stepStatus.Payload.Response = responsePayload
		recordStepCompletionStatus(stepStatus, stepStartTime)
		requestContext.SetStepResponseToContext(step.Name, responsePayload)
		return models.EmptyErrorResponse()
	} else if !step.DidStepExecute() {
		// record step skipped
		// setting response of skipped step with same as request for future validations use
		requestContext.SetStepResponseToContext(step.Name, requestContext.GetStepRequestFromContext(step.Name))
		recordStepSkippedStatus(stepStatus, stepStartTime)
		return models.EmptyErrorResponse()
	}
	return models.EmptyErrorResponse()
}

func recordStepCompletionStatus(stepStatus *models.StepsStatus, stepStartTime time.Time) {
	stepStatus.Status = models.StatusCompleted
	stepStatus.TotalTimeInMs = time.Since(stepStartTime).Nanoseconds() / utils.MilliSecondsDivisor
	SaveStepStatus(stepStatus)
}

func recordStepSkippedStatus(stepStatus *models.StepsStatus, stepStartTime time.Time) {
	stepStatus.Status = models.StatusSkipped
	stepStatus.TotalTimeInMs = time.Since(stepStartTime).Nanoseconds() / utils.MilliSecondsDivisor
	SaveStepStatus(stepStatus)
}

func recordStepPausedStatus(stepStatus *models.StepsStatus, stepStartTime time.Time) {
	stepStatus.Status = models.StatusPaused
	stepStatus.TotalTimeInMs = time.Since(stepStartTime).Nanoseconds() / utils.MilliSecondsDivisor
	SaveStepStatus(stepStatus)
}

func recordStepStartedStatus(stepStatus *models.StepsStatus, stepStartTime time.Time) {
	stepStatus.Status = models.StatusStarted
	stepStatus.TotalTimeInMs = time.Since(stepStartTime).Nanoseconds() / utils.MilliSecondsDivisor
	SaveStepStatus(stepStatus)
}

func recordStepFailedStatus(stepStatus *models.StepsStatus, clampErrorResponse models.ClampErrorResponse, stepStartTime time.Time) {
	stepStatus.Status = models.StatusFailed
	marshal, err := json.Marshal(clampErrorResponse)
	if err != nil {
		log.Error("clampErrorResponse: Marshal error", err)
		return
	}

	var responsePayload map[string]interface{}
	err = json.Unmarshal(marshal, &responsePayload)
	if err != nil {
		log.Error("clampErrorResponse: UnMarshal error", err)
		return
	}

	errPayload := map[string]interface{}{"errors": responsePayload}
	stepStatus.Payload.Response = errPayload
	stepStatus.Reason = clampErrorResponse.Message
	stepStatus.TotalTimeInMs = time.Since(stepStartTime).Nanoseconds() / utils.MilliSecondsDivisor
	SaveStepStatus(stepStatus)
}

func getServiceRequestChannel() chan models.ServiceRequest {
	if serviceRequestChannel == nil {
		panic(errors.New("service request channel not initialized"))
	}
	return serviceRequestChannel
}

func AddServiceRequestToChannel(serviceReq *models.ServiceRequest) {
	channel := getServiceRequestChannel()
	channel <- *serviceReq
}
