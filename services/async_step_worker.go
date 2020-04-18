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
const AsyncStepExecutorWorkersSize = 100

var (
	asyncChannel chan models.AsyncStepRequest
	singletonAsyncStepExecutor        sync.Once
)

func CreateAsyncStepExecutorChannel() chan models.AsyncStepRequest {
	singletonAsyncStepExecutor.Do(func() {
		asyncChannel = make(chan models.AsyncStepRequest, AsyncStepExecutorChannelSize)
	})
	return asyncChannel
}

func init() {
	CreateAsyncStepExecutorChannel()
	CreateAsyncStepExecutorWorkers()
}

func CreateAsyncStepExecutorWorkers() {
	for i := 0; i < AsyncStepExecutorWorkersSize; i++ {
		go asyncWorker(i, asyncChannel)
	}
}

func asyncWorker(workerId int, asyncServiceReqChan <-chan models.AsyncStepRequest) {
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
			AddResumeStepResponseToChannel(resumeRequest)
			return
		}
	}
}

func GetAsyncStepExecutorChannel() chan models.AsyncStepRequest {
	if asyncChannel == nil {
		panic(errors.New("async service request channel not initialized"))
	}
	return asyncChannel
}

func AddAsyncStepExecutorRequestToChannel(asyncStepExecutorRequest models.AsyncStepRequest) {
	channel := GetAsyncStepExecutorChannel()
	channel <- asyncStepExecutorRequest
}
