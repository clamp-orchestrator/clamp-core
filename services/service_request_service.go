package services

import (
	"clamp-core/models"
	"clamp-core/repository"
	"github.com/google/uuid"
	"log"
)

//FindServiceRequestByID is
func FindServiceRequestByID(serviceRequestId uuid.UUID) (models.ServiceRequest, error) {
	log.Printf("Find service Request request by id: %s", serviceRequestId)
	serviceRequest, err := repository.GetDB().FindServiceRequestById(serviceRequestId)
	if err != nil {
		log.Printf("No record found with given service request id %s", serviceRequestId)
	}
	return serviceRequest, err
}

func SaveServiceRequest(serviceReq models.ServiceRequest) (models.ServiceRequest, error) {
	log.Printf("Saving service request: %v", serviceReq)
	serviceRequest, err := repository.GetDB().SaveServiceRequest(serviceReq)
	if err != nil {
		log.Printf("Failed saving service request %v, error: %s", serviceRequest, err.Error())
	}
	return serviceRequest, err
}

func FindServiceRequestByWorkflowName(workflowName string, pageNumber int, pageSize int) ([]models.ServiceRequest, error) {
	log.Printf("Getting service request by workflow name: %s", workflowName)
	serviceRequests, err := repository.GetDB().FindServiceRequestsByWorkflowName(workflowName, pageNumber, pageSize)
	if err != nil {
		log.Printf("Failed to fetch service requests by workflow nam: %s for pageNumber: %d, pageSize: %d", workflowName, pageNumber, pageSize)
	}
	return serviceRequests, err
}
