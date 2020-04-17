package services

import (
	"clamp-core/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"sync"
	"time"
)

const ServiceRequestChannelSize = 1000
const ServiceRequestWorkersSize = 100

var (
	serviceRequestChannel chan models.ServiceRequest
	singletonOnce         sync.Once
)

func createServiceRequestChannel() chan models.ServiceRequest {
	singletonOnce.Do(func() {
		serviceRequestChannel = make(chan models.ServiceRequest, ServiceRequestChannelSize)
	})
	return serviceRequestChannel
}

func init() {
	createServiceRequestChannel()
	createServiceRequestWorkers()
}

func createServiceRequestWorkers() {
	for i := 0; i < ServiceRequestWorkersSize; i++ {
		go worker(i, serviceRequestChannel)
	}
}

func worker(workerId int, serviceReqChan <-chan models.ServiceRequest) {
	prefix := fmt.Sprintf("[WORKER_%d] : ", workerId)
	prefix = fmt.Sprintf("%15s", prefix)
	log.Printf("%s Started listening to service request channel\n", prefix)
	for serviceReq := range serviceReqChan {
		executeWorkflow(serviceReq, prefix)
	}
}

func executeWorkflow(serviceReq models.ServiceRequest, prefix string) {
	prefix = prefix[:len(prefix)-2] + fmt.Sprintf("[REQUEST_ID: %s]", serviceReq.ID)
	log.Printf("%s Started processing service request id %s\n", prefix, serviceReq.ID)
	defer catchErrors(prefix, serviceReq.ID)

	start := time.Now()
	workflow, err := FindWorkflowByName(serviceReq.WorkflowName)
	if err == nil {
		log.Println("Inside Async Execution mode")
		executeWorkflowStepsInSync(workflow, prefix, serviceReq)
	}
	elapsed := time.Since(start)
	log.Printf("%s Completed processing service request id %s in %s\n", prefix, serviceReq.ID, elapsed)
}

func catchErrors(prefix string, requestId uuid.UUID) {
	if r := recover(); r != nil {
		log.Println("[ERROR]", r)
		log.Printf("%s Failed processing service request id %s\n", prefix, requestId)
	}
}

func executeWorkflowStepsInSync(workflow models.Workflow, prefix string, serviceRequest models.ServiceRequest) {
	previousStepResponse := serviceRequest.Payload
	var stepStatus models.StepsStatus
	stepStatus.WorkflowName = serviceRequest.WorkflowName
	stepStatus.ServiceRequestId = serviceRequest.ID

	for _, step := range workflow.Steps {
		if serviceRequest.CurrentStepId != 0 {
			if step.Id == serviceRequest.CurrentStepId{
				serviceRequest.CurrentStepId = 0
			}
			continue
		}
		if step.StepType == "SYNC" {
			ExecuteWorkflowStep(stepStatus, previousStepResponse, prefix, step)
		} else {
			asyncStepExecutionRequest := PrepareAsyncStepExecutionRequest(stepStatus, step, previousStepResponse, prefix)
			AddAsyncStepExecutionRequestToChannel(asyncStepExecutionRequest)
			return
		}
	}
}

func PrepareAsyncStepExecutionRequest(stepStatus models.StepsStatus, step models.Step, previousStepResponse map[string]interface{}, prefix string) models.AsyncStepExecutionRequest {
	asyncStepExecutionRequest := models.AsyncStepExecutionRequest{
		StepStatus: stepStatus,
		Step:       step,
		Payload:    previousStepResponse,
		Prefix:     prefix,
	}
	return asyncStepExecutionRequest
}

func ExecuteWorkflowStep(stepStatus models.StepsStatus, previousStepResponse map[string]interface{}, prefix string, step models.Step) (interface{}, error) {
	stepStatus.Payload.Request = previousStepResponse
	stepStatus.Payload.Response = nil
	stepStartTime := time.Now()
	stepStatus.StepName = step.Name
	recordStepStartedStatus(stepStatus, stepStartTime)
	oldPrefix := log.Prefix()
	log.SetPrefix(oldPrefix + prefix)
	resp, err := step.DoExecute(stepStatus.Payload.Request)
	log.SetPrefix(oldPrefix)
	if err != nil {
		log.Println("Inside error block", err)
		recordStepFailedStatus(stepStatus, err, stepStartTime, prefix)
		errFmt := fmt.Errorf("%s Failed executing step %s, %s \n", prefix, stepStatus.StepName, err.Error())
		panic(errFmt)
	}
	if resp != nil {
		log.Printf("%s Received %s", prefix, resp.(string))
		var responsePayload map[string]interface{}
		json.Unmarshal([]byte(resp.(string)), &responsePayload)
		stepStatus.Payload.Response = responsePayload
		recordStepCompletionStatus(stepStatus, stepStartTime)
		previousStepResponse = responsePayload
	}
	return resp, err
}

func recordStepCompletionStatus(stepStatus models.StepsStatus, stepStartTime time.Time) {
	stepStatus.Status = models.STATUS_COMPLETED
	stepStatus.TotalTimeInMs = time.Since(stepStartTime).Nanoseconds() / 1000000
	SaveStepStatus(stepStatus)
}

func recordStepPausedStatus(stepStatus models.StepsStatus, stepStartTime time.Time) {
	stepStatus.Status = models.STATUS_PAUSED
	stepStatus.TotalTimeInMs = time.Since(stepStartTime).Nanoseconds() / 1000000
	SaveStepStatus(stepStatus)
}

func recordStepStartedStatus(stepStatus models.StepsStatus, stepStartTime time.Time) {
	stepStatus.Status = models.STATUS_STARTED
	stepStatus.TotalTimeInMs = time.Since(stepStartTime).Nanoseconds() / 1000000
	SaveStepStatus(stepStatus)
}

func recordStepFailedStatus(stepStatus models.StepsStatus, err error, stepStartTime time.Time, prefix string) {
	stepStatus.Status = models.STATUS_FAILED
	clampErrorResponse := models.CreateErrorResponse(http.StatusBadRequest, err.Error())
	marshal, marshalErr := json.Marshal(clampErrorResponse)
	log.Println("clampErrorResponse: Marshal error", marshalErr)
	var responsePayload map[string]interface{}
	unmarshalErr := json.Unmarshal(marshal, &responsePayload)
	log.Println("clampErrorResponse: UnMarshal error", unmarshalErr)
	errPayload := map[string]interface{}{"errors": responsePayload}
	stepStatus.Payload.Response = errPayload
	stepStatus.Reason = err.Error()
	stepStatus.TotalTimeInMs = time.Since(stepStartTime).Nanoseconds() / 1000000
	SaveStepStatus(stepStatus)
}

func getServiceRequestChannel() chan models.ServiceRequest {
	if serviceRequestChannel == nil {
		panic(errors.New("service request channel not initialized"))
	}
	return serviceRequestChannel
}

func AddServiceRequestToChannel(serviceReq models.ServiceRequest) {
	channel := getServiceRequestChannel()
	channel <- serviceReq
}
