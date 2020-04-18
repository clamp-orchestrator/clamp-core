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
const ResumeStepResponseChannelSize = 1000
const ResumeStepResponseWorkersSize = 100

var (
	resumeStepResponseChannel chan models.ResumeStepResponse
	singleton                 sync.Once
)

func CreateResumeStepResponseChannel() chan models.ResumeStepResponse {
	singleton.Do(func() {
		resumeStepResponseChannel = make(chan models.ResumeStepResponse, ResumeStepResponseChannelSize)
	})
	return resumeStepResponseChannel
}

func init() {
	CreateResumeStepResponseChannel()
	CreateResumeStepResponseWorkers()
}

func CreateResumeStepResponseWorkers() {
	for i := 0; i < ResumeStepResponseWorkersSize; i++ {
		go resumeStepResponseWorker(i, resumeStepResponseChannel)
	}
}

func resumeStepResponseWorker(workerId int, resumeStepResponsesChan <-chan models.ResumeStepResponse) {
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

func GetResumeStepResponseChannel() chan models.ResumeStepResponse {
	if resumeStepResponseChannel == nil {
		panic(errors.New("async service request channel not initialized"))
	}
	return resumeStepResponseChannel
}

func AddResumeStepResponseToChannel(asyncResumeStepExecutionRequest models.ResumeStepResponse) {
	channel := GetResumeStepResponseChannel()
	channel <- asyncResumeStepExecutionRequest
}
