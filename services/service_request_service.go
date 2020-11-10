package services

import (
	"clamp-core/models"
	"clamp-core/repository"
	"log"

	"github.com/google/uuid"
)

//FindServiceRequestByID is used to fetch service requests by their ID values
func FindServiceRequestByID(serviceRequestID uuid.UUID) (models.ServiceRequest, error) {
	log.Printf("Find service Request request by id: %s", serviceRequestID)
	serviceRequest, err := repository.GetDB().FindServiceRequestByID(serviceRequestID)
	if err != nil {
		log.Printf("No record found with given service request id %s", serviceRequestID)
	}
	return serviceRequest, err
}

//SaveServiceRequest is used to save the created service requests to DB
func SaveServiceRequest(serviceReq models.ServiceRequest) (models.ServiceRequest, error) {
	log.Printf("Saving service request: %v", serviceReq)
	serviceRequest, err := repository.GetDB().SaveServiceRequest(serviceReq)
	if err != nil {
		log.Printf("Failed saving service request %v, error: %s", serviceRequest, err.Error())
	}
	return serviceRequest, err
}

//FindServiceRequestByWorkflowName fetches all ServiceRequests that are associated to a workflow type
func FindServiceRequestByWorkflowName(workflowName string, pageNumber int, pageSize int, sortBy models.SortByFields) ([]models.ServiceRequest, int, error) {
	log.Printf("Getting service request by workflow name: %s", workflowName)
	serviceRequests, totalServiceRequestsCount, err := repository.GetDB().FindServiceRequestsByWorkflowName(workflowName, pageNumber, pageSize, sortBy)
	if err != nil {
		log.Printf("Failed to fetch service requests by workflow nam: %s for pageNumber: %d, pageSize: %d", workflowName, pageNumber, pageSize)
	}
	return serviceRequests, totalServiceRequestsCount, err
}
