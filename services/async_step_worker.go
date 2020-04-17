package services

import (
	"clamp-core/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
)

const AsyncStepExecutionChannelSize = 1000
const AsyncStepExecutionWorkersSize = 100

var (
	asyncChannel chan models.AsyncStepExecutionRequest
	singletonAsyncStepExecutor        sync.Once
)

func CreateAsyncStepExecutionChannel() chan models.AsyncStepExecutionRequest {
	singletonAsyncStepExecutor.Do(func() {
		asyncChannel = make(chan models.AsyncStepExecutionRequest, AsyncStepExecutionChannelSize)
	})
	return asyncChannel
}

func init() {
	CreateAsyncStepExecutionChannel()
	CreateAsyncStepExecutionWorkers()
}

func CreateAsyncStepExecutionWorkers() {
	for i := 0; i < AsyncStepExecutionWorkersSize; i++ {
		go asyncWorker(i, asyncChannel)
	}
}

func asyncWorker(workerId int, asyncServiceReqChan <-chan models.AsyncStepExecutionRequest) {
	prefix := fmt.Sprintf("[WORKER_%d] : ", workerId)
	prefix = fmt.Sprintf("%15s", prefix)
	log.Printf("%s Started listening to service request channel\n", prefix)
	for asyncServiceReq := range asyncServiceReqChan {
		asyncStepResponse, _ := ExecuteWorkflowStep(asyncServiceReq.StepStatus, asyncServiceReq.Payload, asyncServiceReq.Prefix, asyncServiceReq.Step)
		if asyncStepResponse != nil {
			var responsePayload map[string]interface{}
			json.Unmarshal([]byte(asyncStepResponse.(string)), &responsePayload)
			resumeRequest := models.ResumeStepResponse{
				ServiceRequestId: asyncServiceReq.StepStatus.ServiceRequestId,
				StepId:           asyncServiceReq.Step.Id,
				Payload:          responsePayload,
				StepProcessed:    true,
			}
			AddAsyncResumeStepExecutionRequestToChannel(resumeRequest)
			return
		}
	}
}

func GetAsyncStepExecutionChannel() chan models.AsyncStepExecutionRequest {
	if asyncChannel == nil {
		panic(errors.New("async service request channel not initialized"))
	}
	return asyncChannel
}

func AddAsyncStepExecutionRequestToChannel(asyncStepExecutionRequest models.AsyncStepExecutionRequest) {
	channel := GetAsyncStepExecutionChannel()
	channel <- asyncStepExecutionRequest
}
