package services

import (
	"clamp-core/models"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"sync"
)

//TODO Channel name to be changed
const ResumeStepResponseChannelSize = 1000
const ResumeStepResponseWorkersSize = 100

var (
	resumeStepsChannel chan models.AsyncStepResponse
	singleton          sync.Once
)

func createResumeStepsChannel() chan models.AsyncStepResponse {
	singleton.Do(func() {
		resumeStepsChannel = make(chan models.AsyncStepResponse, ResumeStepResponseChannelSize)
	})
	return resumeStepsChannel
}

func init() {
	createResumeStepsChannel()
	createResumeStepsWorkers()
}

func createResumeStepsWorkers() {
	for i := 0; i < ResumeStepResponseWorkersSize; i++ {
		go resumeSteps(i, resumeStepsChannel)
	}
}

func resumeSteps(workerId int, resumeStepsChannel <-chan models.AsyncStepResponse) {
	prefix := fmt.Sprintf("[RESUME_STEP_WORKER_%d] : ", workerId)
	prefix = fmt.Sprintf("%15s", prefix)
	log.Printf("%s Started listening to service request channel\n", prefix)
	for stepResponse := range resumeStepsChannel {
		if !stepResponse.StepProcessed {
			//Fetch from DB the last step executed
			currentStepStatus, _ := FindStepStatusByServiceRequestIdAndStepNameAndStatus(stepResponse.ServiceRequestId, "stepResponse.StepId", models.STATUS_STARTED)
			currentStepStatus.ID = ""
			//TODO Setting Id empty and also errors validations
			if stepResponse.Errors.Code == 0 {
				currentStepStatus.Payload.Response = stepResponse.Payload
				recordStepCompletionStatus(currentStepStatus, currentStepStatus.CreatedAt)
			} else {
				recordStepFailedStatus(currentStepStatus, stepResponse.Errors, currentStepStatus.CreatedAt)
				return
			}
		}
		serviceRequest, err := FindServiceRequestByID(stepResponse.ServiceRequestId)
		if err == nil {
			//TODO
		}
		serviceRequest.Payload = stepResponse.Payload
		serviceRequest.CurrentStepId = stepResponse.StepId
		AddServiceRequestToChannel(serviceRequest)
	}
}

func getResumeStepResponseChannel() chan models.AsyncStepResponse {
	if resumeStepsChannel == nil {
		panic(errors.New("async service request channel not initialized"))
	}
	return resumeStepsChannel
}

func AddStepResponseToResumeChannel(response models.AsyncStepResponse) {
	if response.ServiceRequestId == uuid.Nil || response.StepId == 0 || response.Payload == nil {
		log.Printf("Invalid step resume request received %v", response)
		return
	}
	channel := getResumeStepResponseChannel()
	channel <- response
}
