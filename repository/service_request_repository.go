package repository

import (
	"clamp-core/config"
	"clamp-core/models"

	"github.com/google/uuid"
)

// DBInterface provides a collection of method signatures that needs to be implemented for a specific database.
type ServiceRequestRepository interface {
	SaveServiceRequest(*models.ServiceRequest) (*models.ServiceRequest, error)
	FindServiceRequestByID(uuid.UUID) (*models.ServiceRequest, error)
	FindServiceRequestsByWorkflowName(workflowName string, pageNumber int, pageSize int) ([]*models.ServiceRequest, error)
}

var serviceRequestRepository ServiceRequestRepository

func init() {
	switch config.ENV.DBDriver {
	case "postgres":
		serviceRequestRepository = &servicerequestrepositorypostgres{}
	}
}

// GetDB returns the initialized database implementations. Currently only postgres is implemented.
func GetServiceRequestRepository() ServiceRequestRepository {
	return serviceRequestRepository
}

// SetDB is used to update the db object with custom implementations.
// It is used in tests to override the actual db implementations with mock implementations
func SetServiceRequestRepository(serviceRequestRepositoryImpl ServiceRequestRepository) {
	serviceRequestRepository = serviceRequestRepositoryImpl
}
