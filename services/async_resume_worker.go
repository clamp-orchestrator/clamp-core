package services

import (
	"clamp-core/models"
	"errors"
	"fmt"
	"log"
	"sync"
)
//TODO Channel name to be changed
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

func asyncResumeWorker(workerId int, asyncResumeReqChan <-chan models.AsyncResumeStepExecutionRequest) {
	prefix := fmt.Sprintf("[WORKER_%d] : ", workerId)
	prefix = fmt.Sprintf("%15s", prefix)
	log.Printf("%s Started listening to service request channel\n", prefix)
	for resumeRequest := range asyncResumeReqChan {
		log.Println("------------- Resume Request --------------", resumeRequest)
		//TODO Step Processed will be set to false by default, if async http has response then it will be set to true,
		if !resumeRequest.StepProcessed {

		}
		//if true then skip marking that step as Completed
		// Instead directly call next step to execute
		// Check if payload contains error block
		// if so then mark step as failed and stop.....
		serviceRequest, err := FindServiceRequestByID(resumeRequest.ServiceRequestId)
		if err == nil {
			// TODO Handle error case
		}
		serviceRequest.Payload = resumeRequest.Payload
		serviceRequest.CurrentStepId = resumeRequest.StepId
		AddServiceRequestToChannel(serviceRequest)
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
