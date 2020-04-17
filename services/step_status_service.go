package services

import (
	"clamp-core/models"
	"clamp-core/repository"
	"github.com/google/uuid"
	"log"
	"time"
)

func SaveStepStatus(stepStatusReq models.StepsStatus) (models.StepsStatus, error) {
	log.Printf("Saving step status : %v", stepStatusReq)
	stepStatusReq, err := repository.GetDB().SaveStepStatus(stepStatusReq)
	if err != nil {
		log.Printf("Failed saving step status : %v, %s", stepStatusReq, err.Error())
	}
	return stepStatusReq, err
}

func FindStepStatusByServiceRequestId(serviceRequestId uuid.UUID) ([]models.StepsStatus, error) {
	log.Printf("Find step statues by request id : %s ", serviceRequestId)
	stepsStatuses, err := repository.GetDB().FindStepStatusByServiceRequestId(serviceRequestId)
	if err != nil {
		log.Printf("No record found with given service request id %s", serviceRequestId)
		return []models.StepsStatus{}, err
	}
	return stepsStatuses, err
}

func FindStepStatusByServiceRequestIdAndStatusOrderByCreatedAtDesc(serviceRequestId uuid.UUID, status models.Status) (models.StepsStatus, error) {
	log.Printf("Find step statues by request id : %s ", serviceRequestId)
	stepsStatuses, err := repository.GetDB().FindStepStatusByServiceRequestIdAndStatusOrderByCreatedAtDesc(serviceRequestId, status)
	if err != nil {
		log.Printf("No record found with given service request id %s", serviceRequestId)
		return models.StepsStatus{}, err
	}
	return stepsStatuses, err
}

func PrepareStepStatusResponse(stepsStatusArr []models.StepsStatus) models.StepsStatusResponse {
	var stepsStatusRes models.StepsStatusResponse
	steps := make([]models.StepResponse, len(stepsStatusArr))
	if len(stepsStatusArr) > 0 {
		var statusFlag = true

		for i, stepsStatus := range stepsStatusArr {
			stepsStatusRes.Reason = stepsStatus.Reason
			if models.STATUS_FAILED == stepsStatus.Status {
				statusFlag = false
				stepsStatusRes.Status = stepsStatus.Status
			}
			steps[i] = models.StepResponse{
				Name:      stepsStatus.StepName,
				Status:    stepsStatus.Status,
				TimeTaken: stepsStatus.TotalTimeInMs,
				Payload:   stepsStatus.Payload,
			}
		}

		if statusFlag && len(stepsStatusArr)/2 == 0{
			stepsStatusRes.Status = models.STATUS_COMPLETED
		}else if !statusFlag{
			stepsStatusRes.Status = models.STATUS_FAILED
		} else {
			stepsStatusRes.Status = models.STATUS_PAUSED
		}
		stepsStatusRes.ServiceRequestId = stepsStatusArr[0].ServiceRequestId
		stepsStatusRes.WorkflowName = stepsStatusArr[0].WorkflowName
		timeTaken := calculateTimeTaken(stepsStatusArr[0].CreatedAt, stepsStatusArr[len(stepsStatusArr)-1].CreatedAt)
		stepsStatusRes.TotalTimeInMs = timeTaken.Nanoseconds() / 1000000
		stepsStatusRes.Steps = steps
	}
	return stepsStatusRes
}

func calculateTimeTaken(startTime time.Time, endTime time.Time) time.Duration {
	log.Println("Time Difference is == ", endTime.Sub(startTime))
	return endTime.Sub(startTime)
}
