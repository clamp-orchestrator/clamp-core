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
	prefix := fmt.Sprintf("[WORKER_%d] ", workerId)
	prefix = fmt.Sprintf("%15s", prefix)
	log.Printf("%s : Started listening to service request channel\n", prefix)
	for serviceReq := range serviceReqChan {
		executeWorkflow(serviceReq, prefix)
	}
}

func executeWorkflow(serviceReq models.ServiceRequest, prefix string) {
	prefix = fmt.Sprintf("%s [REQUEST_ID: %s]", prefix, serviceReq.ID)
	log.Printf("%s Started processing service request id %s\n", prefix, serviceReq.ID)
	start := time.Now()
	workflow, err := FindWorkflowByName(serviceReq.WorkflowName)
	if err == nil {
		lastStep := workflow.Steps[len(workflow.Steps)-1]
		if serviceReq.CurrentStepId == 0 || serviceReq.CurrentStepId != lastStep.Id {
			executeWorkflowStepsInSync(workflow, prefix, serviceReq)
		} else {
			log.Printf("%s All steps are executed for service request id: %s\n", prefix, serviceReq.ID)
		}
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
	stepRequestPayload := serviceRequest.Payload
	lastStepExecuted := serviceRequest.CurrentStepId
	executeStepsFromIndex := 0
	if lastStepExecuted > 0 {
		executeStepsFromIndex = lastStepExecuted
		log.Printf("%s Skipping steps till  step id %d\n", prefix, executeStepsFromIndex)
	}
	var requestContext models.RequestContext
	requestContext.ServiceRequestId = serviceRequest.ID
	requestContext.WorkflowName = serviceRequest.WorkflowName
	stepsRequestResponsePayload := make(map[string]models.RequestResponse)
	serviceRequest.RequestContext = requestContext
	for _, step := range workflow.Steps[executeStepsFromIndex:] {
		if step.StepType == "SYNC" {
			stepResponsePayload, _ := ExecuteWorkflowStep(stepRequestPayload, serviceRequest.ID, serviceRequest.WorkflowName, step, prefix)
			UpdateStepRequestResponseInRequestContext(step, serviceRequest, stepResponsePayload, stepsRequestResponsePayload, requestContext)
		} else {
			asyncStepExecutionRequest := prepareAsyncStepExecutionRequest(step, stepRequestPayload, serviceRequest)
			UpdateStepRequestResponseInRequestContext(step, serviceRequest, nil, stepsRequestResponsePayload, requestContext)
			AddAsyncStepToExecutorChannel(asyncStepExecutionRequest)
			return
		}
		requestContext.Payload = stepsRequestResponsePayload
		log.Println(" ----========= Request Context Object Payload ----=========", requestContext.Payload)
		serviceRequest.RequestContext = requestContext
	}
}

func UpdateStepRequestResponseInRequestContext(step models.Step, serviceRequest models.ServiceRequest, stepResponsePayload map[string]interface{}, requestResponsePayload map[string]models.RequestResponse, requestContext models.RequestContext) {
	stepRequestResponse := map[string]models.RequestResponse{step.Name: {
		Request:  serviceRequest.Payload,
		Response: stepResponsePayload,
	}}
	requestResponsePayload[step.Name] = stepRequestResponse[step.Name]
}

func prepareAsyncStepExecutionRequest(step models.Step, previousStepResponse map[string]interface{}, serviceReq models.ServiceRequest) models.AsyncStepRequest {
	asyncStepExecutionRequest := models.AsyncStepRequest{
		Step:             step,
		Payload:          previousStepResponse,
		ServiceRequestId: serviceReq.ID,
		WorkflowName:     serviceReq.WorkflowName,
	}
	return asyncStepExecutionRequest
}

//TODO: replace prefix with other standard way like MDC
func ExecuteWorkflowStep(stepRequestPayload map[string]interface{}, serviceRequestId uuid.UUID, workflowName string, step models.Step, prefix string) (map[string]interface{}, models.ClampErrorResponse) {
	defer catchErrors(prefix, serviceRequestId)
	stepStartTime := time.Now()
	stepStatus := models.StepsStatus{
		ServiceRequestId: serviceRequestId,
		WorkflowName:     workflowName,
		StepName:         step.Name,
		Payload: models.Payload{
			Request:  stepRequestPayload,
			Response: nil,
		},
		StepId: step.Id,
	}
	recordStepStartedStatus(stepStatus, stepStartTime)
	request := models.StepRequest{
		ServiceRequestId: stepStatus.ServiceRequestId,
		StepId:           step.Id,
		Payload:          stepStatus.Payload.Request,
	}
	resp, err := step.DoExecute(request, prefix)
	if step.DidStepExecute() {
		if err != nil {
			clampErrorResponse := models.CreateErrorResponse(http.StatusBadRequest, err.Error())
			recordStepFailedStatus(stepStatus, *clampErrorResponse, stepStartTime)
			errFmt := fmt.Errorf("%s Failed executing step %s, %s \n", prefix, stepStatus.StepName, err.Error())
			panic(errFmt)
			return nil, *clampErrorResponse
		}
		if resp != nil {
			log.Printf("%s Step response received: %s", prefix, resp.(string))
			var responsePayload map[string]interface{}
			json.Unmarshal([]byte(resp.(string)), &responsePayload)
			stepStatus.Payload.Response = responsePayload
			recordStepCompletionStatus(stepStatus, stepStartTime)
			return responsePayload, models.EmptyErrorResponse()
		}
	} else {
		//record step skipped
		return stepRequestPayload, models.EmptyErrorResponse()
	}
	return nil, models.EmptyErrorResponse()
}

func recordStepCompletionStatus(stepStatus models.StepsStatus, stepStartTime time.Time) {
	stepStatus.Status = models.STATUS_COMPLETED
	stepStatus.TotalTimeInMs = time.Since(stepStartTime).Nanoseconds() / models.MilliSecondsDivisor
	SaveStepStatus(stepStatus)
}

func recordStepPausedStatus(stepStatus models.StepsStatus, stepStartTime time.Time) {
	stepStatus.Status = models.STATUS_PAUSED
	stepStatus.TotalTimeInMs = time.Since(stepStartTime).Nanoseconds() / models.MilliSecondsDivisor
	SaveStepStatus(stepStatus)
}

func recordStepStartedStatus(stepStatus models.StepsStatus, stepStartTime time.Time) {
	stepStatus.Status = models.STATUS_STARTED
	stepStatus.TotalTimeInMs = time.Since(stepStartTime).Nanoseconds() / models.MilliSecondsDivisor
	SaveStepStatus(stepStatus)
}

func recordStepFailedStatus(stepStatus models.StepsStatus, clampErrorResponse models.ClampErrorResponse, stepStartTime time.Time) {
	stepStatus.Status = models.STATUS_FAILED
	marshal, marshalErr := json.Marshal(clampErrorResponse)
	log.Println("clampErrorResponse: Marshal error", marshalErr)
	var responsePayload map[string]interface{}
	unmarshalErr := json.Unmarshal(marshal, &responsePayload)
	log.Println("clampErrorResponse: UnMarshal error", unmarshalErr)
	errPayload := map[string]interface{}{"errors": responsePayload}
	stepStatus.Payload.Response = errPayload
	stepStatus.Reason = clampErrorResponse.Message
	stepStatus.TotalTimeInMs = time.Since(stepStartTime).Nanoseconds() / models.MilliSecondsDivisor
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
