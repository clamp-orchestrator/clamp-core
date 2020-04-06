package services

import (
	"clamp-core/models"
	"clamp-core/repository"
	"fmt"
	"github.com/google/uuid"
	"log"
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
	serviceReq := new(models.StepsStatus)
	var stepStatusArr []models.StepsStatus
	_, err := repository.GetDB().Query(&stepStatusArr, "select * from steps_status where service_request_id = ?", serviceRequestId)
	//err := repo.whereQuery(&stepStatusArr, "steps_status.service_request_id = ?", serviceRequestId)

	if err != nil {
		fmt.Errorf("No record found with given service request id %s", serviceRequestId)
	}

	stepStatusRes := PrepareStepStatusResponse(stepStatusArr)
	log.Println("Steps Status Response is ", stepStatusRes)
	log.Println("Service request id is ", serviceReq.ServiceRequestId)
	return stepStatusRes, err
}

func PrepareStepStatusResponse(stepsStatusArr []models.StepsStatus) models.StepsStatusResponse {
	var stepsStatusRes models.StepsStatusResponse
	steps := make([]models.StepResponse, len(stepsStatusArr))
	var statusFlag bool = true
	for i := range stepsStatusArr{
		stepsStatus := models.StepsStatus{}
		stepsStatus = stepsStatusArr[i]
		stepsStatusRes.ServiceRequestId = stepsStatus.ServiceRequestId
		stepsStatusRes.WorkflowName = stepsStatus.WorkflowName
		stepsStatusRes.Reason = stepsStatus.Reason
		if(stepsStatus.Status == models.STATUS_FAILED){
			statusFlag = false
			stepsStatusRes.Status = stepsStatus.Status
		}
		steps[i] = models.StepResponse{
			Name:      stepsStatus.StepName,
			Status:    stepsStatus.Status,
			TimeTaken: 0,
		}
	}
	if(statusFlag){
		stepsStatusRes.Status = models.STATUS_COMPLETED
	}
	stepsStatusRes.TotalTime = 10
	stepsStatusRes.Steps = steps
	return stepsStatusRes
}