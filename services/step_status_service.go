package services

import (
	"clamp-core/models"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
)

func SaveStepStatus(stepStatusReq models.StepsStatus) (models.StepsStatus, error) {
	pgStepStatusReq := stepStatusReq.ToPgStepStatus()
	err := repo.insertQuery(&pgStepStatusReq)

	if err != nil {
		panic(err)
	}
	return stepStatusReq, err
}

func FindStepStatusByServiceRequestId(serviceRequestId uuid.UUID) (models.StepsStatusResponse, error) {
	serviceRequestReq := models.StepsStatus{ServiceRequestId: serviceRequestId}
	fmt.Println("Service Request request is -- ", serviceRequestReq)
	var stepsStatus []models.StepsStatus

	_, err := repo.query(&stepsStatus, "select * from steps_status where service_request_id = ?", serviceRequestId)
	//err := repo.whereQuery(stepsStatus, "steps_status.service_request_id = ?", serviceRequestId)

	log.Println("Steps Status Where Query Response is ", stepsStatus)
	if err != nil {
		fmt.Errorf("No record found with given service request id %s", serviceRequestId)
		return models.StepsStatusResponse{}, err
	}

	stepStatusRes := PrepareStepStatusResponse(stepsStatus)
	//stepStatusRes := models.StepsStatusResponse{}

	log.Println("Steps Status Response is ", stepStatusRes)
	log.Println("Service request id is ", serviceRequestReq.ServiceRequestId)
	return stepStatusRes, err
}

func PrepareStepStatusResponse(stepsStatusArr []models.StepsStatus) models.StepsStatusResponse {
	var stepsStatusRes models.StepsStatusResponse
	steps := make([]models.StepResponse, len(stepsStatusArr))
	var statusFlag = true

	for i := range stepsStatusArr {
		stepsStatus := models.StepsStatus{}
		stepsStatus = stepsStatusArr[i]
		stepsStatusRes.Reason = stepsStatus.Reason
		if stepsStatus.Status == models.STATUS_FAILED {
			statusFlag = false
			stepsStatusRes.Status = stepsStatus.Status
		}
		steps[i] = models.StepResponse{
			Name:      stepsStatus.StepName,
			Status:    stepsStatus.Status,
			TimeTaken: stepsStatus.TotalTimeInMs,
		}
	}

	if statusFlag {
		stepsStatusRes.Status = models.STATUS_COMPLETED
	}
	stepsStatusRes.ServiceRequestId = stepsStatusArr[0].ServiceRequestId
	stepsStatusRes.WorkflowName = stepsStatusArr[0].WorkflowName
	timeTaken := calculateTimeTaken(stepsStatusArr[0].CreatedAt, stepsStatusArr[len(stepsStatusArr)-1].CreatedAt)
	stepsStatusRes.TotalTimeInMs = timeTaken.Nanoseconds() / 1000
	stepsStatusRes.Steps = steps
	return stepsStatusRes
}

func calculateTimeTaken(startTime time.Time, endTime time.Time) time.Duration {
	log.Println("Time Difference is == ", endTime.Sub(startTime))
	return endTime.Sub(startTime)
}
