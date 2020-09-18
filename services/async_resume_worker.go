package services

import (
	"clamp-core/models"
	"clamp-core/utils"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"sync"
	"time"
)

var (
	resumeStepsChannel chan models.AsyncStepResponse
	singleton          sync.Once
)

func createResumeStepsChannel() chan models.AsyncStepResponse {
	singleton.Do(func() {
		resumeStepsChannel = make(chan models.AsyncStepResponse, utils.ResumeStepResponseChannelSize)
	})
	return resumeStepsChannel
}

func init() {
	createResumeStepsChannel()
	createResumeStepsWorkers()
}

func createResumeStepsWorkers() {
	for i := 0; i < utils.ResumeStepResponseWorkersSize; i++ {
		go resumeSteps(i, resumeStepsChannel)
	}
}

func resumeSteps(workerId int, resumeStepsChannel <-chan models.AsyncStepResponse) {
	duplicateStepResponse := false
	prefix := fmt.Sprintf("[RESUME_STEP_WORKER_%d] ", workerId)
	prefix = fmt.Sprintf("%15s", prefix)
	log.Printf("%s : Started listening to resume steps channel\n", prefix)
	for stepResponse := range resumeStepsChannel {
		prefix = fmt.Sprintf("%s [REQUEST_ID: %s]", prefix, stepResponse.ServiceRequestId)
		log.Printf("%s : Received step response : %v \n", prefix, stepResponse)
		currentStepStatusArr, _ := FindAllStepStatusByServiceRequestIdAndStepId(stepResponse.ServiceRequestId, stepResponse.StepId)
		var currentStepStatus models.StepsStatus
		for _, stepStatus := range currentStepStatusArr {
			if stepStatus.Status == models.STATUS_STARTED {
				currentStepStatus = stepStatus
			}
			if stepStatus.Status == models.STATUS_COMPLETED || stepStatus.Status == models.STATUS_FAILED {
				log.Printf("%s : [DUPLICATE_STEP_RESPONSE] Received step is already completed : %v \n", prefix, stepResponse)
				duplicateStepResponse = true
				break
			}
		}
		if !duplicateStepResponse {
			if !stepResponse.IsStepStatusRecorded() {
				currentStepStatus.ID = ""
				resumeStepStartTime := currentStepStatus.CreatedAt
				currentStepStatus.CreatedAt = time.Time{}
				//TODO Setting Id empty and also errors validations
				//TODO subtracting -5.30 since time is stored in GMT in PSql
				if !stepResponse.Error.IsNil() {
					recordStepFailedStatus(currentStepStatus, stepResponse.Error, currentStepStatus.CreatedAt.Add(-(time.Minute * 330)))
					continue
				} else {
					currentStepStatus.Payload.Response = stepResponse.Response
					recordStepCompletionStatus(currentStepStatus, resumeStepStartTime)
				}
			}
			serviceRequest, err := FindServiceRequestByID(stepResponse.ServiceRequestId)
			if err != nil {
				//TODO
			} else {
				serviceRequest.Payload = stepResponse.Response
				serviceRequest.CurrentStepId = stepResponse.StepId
				serviceRequest.RequestHeaders = stepResponse.RequestHeaders
				AddServiceRequestToChannel(serviceRequest)
			}
		}

	}
}

func getResumeStepResponseChannel() chan models.AsyncStepResponse {
	if resumeStepsChannel == nil {
		panic(errors.New("async service request channel not initialized"))
	}
	return resumeStepsChannel
}

func AddStepResponseToResumeChannel(response models.AsyncStepResponse) {
	if response.ServiceRequestId == uuid.Nil || response.StepId == 0 || (response.Response == nil && response.Error.Code == 0) {
		log.Printf("Invalid step resume request received %v", response)
		return
	}
	channel := getResumeStepResponseChannel()
	channel <- response
}
