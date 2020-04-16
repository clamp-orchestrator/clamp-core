package services

import (
	"clamp-core/models"
	"errors"
	"fmt"
	"log"
	"sync"
)

const AsyncResumeStepExecutionChannelSize = 1000
const AsyncResumeStepExecutionWorkersSize = 100

var (
	asyncResumeStepExecutionChannel chan models.AsyncResumeStepExecutionRequest
	singleton        sync.Once
)

func CreateAsyncResumeStepExecutionChannel() chan models.AsyncResumeStepExecutionRequest {
	singleton.Do(func() {
		asyncResumeStepExecutionChannel = make(chan models.AsyncResumeStepExecutionRequest, AsyncResumeStepExecutionChannelSize)
	})
	return asyncResumeStepExecutionChannel
}

func init() {
	CreateAsyncResumeStepExecutionChannel()
	CreateAsyncResumeStepExecutionWorkers()
}

func CreateAsyncResumeStepExecutionWorkers() {
	for i := 0; i < AsyncResumeStepExecutionWorkersSize; i++ {
		go asyncResumeWorker(i, asyncResumeStepExecutionChannel)
	}
}

func asyncResumeWorker(workerId int, asyncServiceReqChan <-chan models.AsyncResumeStepExecutionRequest) {
	prefix := fmt.Sprintf("[WORKER_%d] : ", workerId)
	prefix = fmt.Sprintf("%15s", prefix)
	log.Printf("%s Started listening to service request channel\n", prefix)
	for asyncServiceReq := range asyncServiceReqChan {
		log.Println("Async Service Request -", asyncServiceReq)
		//ExecuteWorkflowStep(asyncServiceReq.StepStatus,asyncServiceReq.Payload,asyncServiceReq.Prefix, asyncServiceReq.Step)
	}
}

func GetAsyncResumeStepExecutionChannel() chan models.AsyncResumeStepExecutionRequest {
	if asyncResumeStepExecutionChannel == nil {
		panic(errors.New("async service request channel not initialized"))
	}
	return asyncResumeStepExecutionChannel
}

func AddAsyncResumeStepExecutionRequestToChannel(asyncResumeStepExecutionRequest models.AsyncResumeStepExecutionRequest) {
	channel := GetAsyncResumeStepExecutionChannel()
	channel <- asyncResumeStepExecutionRequest
}
