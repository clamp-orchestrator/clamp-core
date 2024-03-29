package services

import (
	"clamp-core/models"
	"clamp-core/utils"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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

func resumeSteps(workerID int, resumeStepsChannel <-chan models.AsyncStepResponse) {
	duplicateStepResponse := false
	prefix := fmt.Sprintf("[RESUME_STEP_WORKER_%d] ", workerID)
	prefix = fmt.Sprintf("%15s", prefix)
	log.Infof("%s : Started listening to resume steps channel", prefix)
	for stepResponse := range resumeStepsChannel {
		prefix = fmt.Sprintf("%s [REQUEST_ID: %s]", prefix, stepResponse.ServiceRequestID)
		log.Debugf("%s : Received step response : %v", prefix, stepResponse)
		currentStepStatusArr, _ := FindAllStepStatusByServiceRequestIDAndStepID(stepResponse.ServiceRequestID, stepResponse.StepID)
		var currentStepStatus *models.StepsStatus
		for _, stepStatus := range currentStepStatusArr {
			if stepStatus.Status == models.StatusStarted {
				currentStepStatus = stepStatus
			}
			if stepStatus.Status == models.StatusCompleted || stepStatus.Status == models.StatusFailed {
				log.Warnf("%s : [DUPLICATE_STEP_RESPONSE] Received step is already completed : %v", prefix, stepResponse)
				duplicateStepResponse = true
				break
			}
		}
		if !duplicateStepResponse {
			if !stepResponse.IsStepStatusRecorded() {
				currentStepStatus.ID = ""
				resumeStepStartTime := currentStepStatus.CreatedAt
				currentStepStatus.CreatedAt = time.Time{}
				// TODO Setting ID empty and also errors validations
				// TODO subtracting -5.30 since time is stored in GMT in PSql
				if !stepResponse.Error.IsNil() {
					recordStepFailedStatus(currentStepStatus, stepResponse.Error, currentStepStatus.CreatedAt.Add(-(time.Minute * 330)))
					continue
				} else {
					currentStepStatus.Payload.Response = stepResponse.Response
					recordStepCompletionStatus(currentStepStatus, resumeStepStartTime)
				}
			}
			serviceRequest, err := FindServiceRequestByID(stepResponse.ServiceRequestID)
			if err != nil {
				// TODO
			} else {
				serviceRequest.Payload = stepResponse.Response
				serviceRequest.CurrentStepID = stepResponse.StepID
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

func AddStepResponseToResumeChannel(response *models.AsyncStepResponse) {
	if response.ServiceRequestID == uuid.Nil || response.StepID == 0 || (response.Response == nil && response.Error.Code == 0) {
		log.Errorf("Invalid step resume request received %v", response)
		return
	}
	channel := getResumeStepResponseChannel()
	channel <- *response
}
