package services

import (
	"clamp-core/models"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)
//TODO Channel name to be changed
const AsyncResumeStepExecutionChannelSize = 1000
const AsyncResumeStepExecutionWorkersSize = 100

var (
	asyncResumeStepExecutionChannel chan models.ResumeStepResponse
	singleton        sync.Once
)

func CreateAsyncResumeStepExecutionChannel() chan models.ResumeStepResponse {
	singleton.Do(func() {
		asyncResumeStepExecutionChannel = make(chan models.ResumeStepResponse, AsyncResumeStepExecutionChannelSize)
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

func asyncResumeWorker(workerId int, resumeStepResponsesChan <-chan models.ResumeStepResponse) {
	prefix := fmt.Sprintf("[WORKER_%d] : ", workerId)
	prefix = fmt.Sprintf("%15s", prefix)
	log.Printf("%s Started listening to service request channel\n", prefix)
	for resumeStepResponse := range resumeStepResponsesChan {
		stepStartTime := time.Now()
		if !resumeStepResponse.StepProcessed {
			//Fetch from DB the last step executed
			currentStepStatus, _ := FindStepStatusByServiceRequestIdAndStatusOrderByCreatedAtDesc(resumeStepResponse.ServiceRequestId, models.STATUS_STARTED)
			currentStepStatus.ID = ""
			//TODO Setting Id empty and also errors validations
			if resumeStepResponse.Errors.Code == 0 {
				currentStepStatus.Payload.Response = resumeStepResponse.Payload
				recordStepCompletionStatus(currentStepStatus, stepStartTime)
			}else{
				recordStepFailedStatus(currentStepStatus,resumeStepResponse.Errors,stepStartTime)
				return
			}
		}
		serviceRequest, err := FindServiceRequestByID(resumeStepResponse.ServiceRequestId)
		if err == nil {
			//TODO
		}
		serviceRequest.Payload = resumeStepResponse.Payload
		serviceRequest.CurrentStepId = resumeStepResponse.StepId
		AddServiceRequestToChannel(serviceRequest)
	}
}

func GetAsyncResumeStepExecutionChannel() chan models.ResumeStepResponse {
	if asyncResumeStepExecutionChannel == nil {
		panic(errors.New("async service request channel not initialized"))
	}
	return asyncResumeStepExecutionChannel
}

func AddAsyncResumeStepExecutionRequestToChannel(asyncResumeStepExecutionRequest models.ResumeStepResponse) {
	channel := GetAsyncResumeStepExecutionChannel()
	channel <- asyncResumeStepExecutionRequest
}
