package services

import (
	"clamp-core/models"
	"github.com/google/uuid"
	"log"
)

//FindServiceRequestByID is
func FindServiceRequestByID(serviceRequestId uuid.UUID) (*models.ServiceRequest, error) {
	serviceRequestReq := models.ServiceRequest{ID: serviceRequestId}
	log.Println("Service Request request is -- ", serviceRequestReq)
	serviceReq := new(models.ServiceRequest)
	err := repo.whereQuery(serviceReq, "service_request.id = ?", serviceRequestId)
	if err != nil {
		log.Printf("No record found with given service request id %s", serviceRequestId)
	}
	return serviceReq, err
}

func SaveServiceRequest(serviceReq models.ServiceRequest) (models.ServiceRequest, error) {
	pgServReq := serviceReq.ToPgServiceRequest()
	err := repo.insertQuery(&pgServReq)
	if err != nil {
		log.Printf("Failed saving service request %v, error: %s", pgServReq, err.Error())
	} else {
		log.Printf("Created new service request %v", pgServReq)
	}
	return pgServReq.ToServiceRequest(), err
}
