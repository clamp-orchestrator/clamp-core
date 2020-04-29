package services

import (
	"clamp-core/models"
	"errors"
	"fmt"
	"log"
	"sync"
)

const AsyncStepExecutorChannelSize = 1000
const AsyncStepExecutorsSize = 100

var (
	asyncStepExecutorChannel   chan models.AsyncStepRequest
	singletonAsyncStepExecutor sync.Once
)

func createAsyncStepExecutorChannel() chan models.AsyncStepRequest {
	singletonAsyncStepExecutor.Do(func() {
		asyncStepExecutorChannel = make(chan models.AsyncStepRequest, AsyncStepExecutorChannelSize)
	})
	return asyncStepExecutorChannel
}

func init() {
	createAsyncStepExecutorChannel()
	createAsyncStepExecutors()
}

func createAsyncStepExecutors() {
	for i := 0; i < AsyncStepExecutorsSize; i++ {
		go executeStep(i, asyncStepExecutorChannel)
	}
}

func executeStep(workerId int, asyncStepExecutorChannel <-chan models.AsyncStepRequest) {
	prefix := fmt.Sprintf("[ASYNC_STEP_EXECUTOR_%d] ", workerId)
	prefix = fmt.Sprintf("%15s", prefix)
	log.Printf("%s : Started listening to async steps executor channel\n", prefix)
	for asyncStepReq := range asyncStepExecutorChannel {
		prefix = fmt.Sprintf("%s [REQUEST_ID: %s]", prefix, asyncStepReq.ServiceRequestId)
		log.Printf("%s : Received async step to execute %v\n", prefix, asyncStepReq)
		workflow, _ := FindWorkflowByName(asyncStepReq.WorkflowName)
		serviceRequest, _ := FindServiceRequestByID(asyncStepReq.ServiceRequestId)
		context := CreateRequestContext(workflow, serviceRequest)
		err := ExecuteWorkflowStep(asyncStepReq.Step, context, prefix)
		if !err.IsNil() {
			asyncStepRes := models.AsyncStepResponse{
				ServiceRequestId: asyncStepReq.ServiceRequestId,
				StepId:           asyncStepReq.Step.Id,
				Response:         nil,
				Error:            err,
			}
			asyncStepRes.SetStepStatusRecorded(true)
			AddStepResponseToResumeChannel(asyncStepRes)
		}
	}
}

func getAsyncStepExecutorChannel() chan models.AsyncStepRequest {
	if asyncStepExecutorChannel == nil {
		panic(errors.New("async service request channel not initialized"))
	}
	return asyncStepExecutorChannel
}

func AddAsyncStepToExecutorChannel(asyncStepRequest models.AsyncStepRequest) {
	channel := getAsyncStepExecutorChannel()
	channel <- asyncStepRequest
}
