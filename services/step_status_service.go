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


func FindStepStatusByServiceRequestId(serviceRequestId uuid.UUID) ([]models.StepsStatus, error) {
	serviceRequestReq := models.StepsStatus{ServiceRequestId: serviceRequestId}
	fmt.Println("Service Request request is -- ", serviceRequestReq)
	serviceReq := new(models.StepsStatus)
	var stepStatusArr []models.StepsStatus
	_, err := repository.GetDB().Query(&stepStatusArr, "select * from steps_status where service_request_id = ?", serviceRequestId)
	//err := repo.whereQuery(&stepStatusArr, "steps_status.service_request_id = ?", serviceRequestId)

	//_,err := repo.query(&stepStatusArr,"select * from steps_status where service_request_id = ?", serviceRequestId)
	if err != nil {
		fmt.Errorf("No record found with given service request id %s", serviceRequestId)
	}
	//return serviceReq, err
	log.Println("Response is ", stepStatusArr)
	log.Println("Service request id is ", serviceReq.ServiceRequestId)
	return stepStatusArr, err
}