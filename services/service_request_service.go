package services

import (
	"clamp-core/models"
	"fmt"
	"github.com/google/uuid"
)

//FindServiceRequestByID is
func FindServiceRequestByID(serviceRequestId uuid.UUID)(*models.ServiceRequest, error) {
	serviceRequestReq := models.ServiceRequest{ID: serviceRequestId}
	fmt.Println("Service Request request is -- ", serviceRequestReq)
	serviceReq := new(models.ServiceRequest)
	err := repo.whereQuery(serviceReq, "service_request.id = ?", serviceRequestId)
	if err != nil {
		fmt.Errorf("No record found with given service request id %s", serviceRequestId)
	}
	return serviceReq, err
}

func SaveServiceRequest(serviceReq models.ServiceRequest) (models.ServiceRequest, error) {
	pgServReq := serviceReq.ToPgServiceRequest()
	err := repo.insertQuery(&pgServReq)

	if err != nil {
		panic(err)
	}
	return serviceReq, err
}
