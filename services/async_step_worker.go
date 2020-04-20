package services

import (
	"clamp-core/models"
	"encoding/json"
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
	prefix := fmt.Sprintf("[ASYNC_STEP_EXECUTOR_%d] : ", workerId)
	prefix = fmt.Sprintf("%15s", prefix)
	log.Printf("%s Started listening to service request channel\n", prefix)
	for asyncStepReq := range asyncStepExecutorChannel {
		asyncStepResponse, _ := ExecuteWorkflowStep(asyncStepReq.StepStatus, asyncStepReq.Payload, asyncStepReq.Prefix, asyncStepReq.Step)
		if asyncStepResponse != nil {
			var responsePayload map[string]interface{}
			json.Unmarshal([]byte(asyncStepResponse.(string)), &responsePayload)
			asyncStepRes := models.AsyncStepResponse{
				ServiceRequestId: asyncStepReq.StepStatus.ServiceRequestId,
				StepId:           asyncStepReq.Step.Id,
				Payload:          responsePayload,
				StepProcessed:    true,
			}
			AddStepResponseToResumeChannel(asyncStepRes)
			return
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
