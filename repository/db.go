package repository

import (
	"clamp-core/config"
	"clamp-core/models"

	"github.com/google/uuid"
)

// DBInterface provides a collection of method signatures that needs to be implemented for a specific database.
type DBInterface interface {
	SaveServiceRequest(models.ServiceRequest) (models.ServiceRequest, error)
	FindServiceRequestByID(uuid.UUID) (models.ServiceRequest, error)
	SaveWorkflow(models.Workflow) (models.Workflow, error)
	FindWorkflowByName(string) (models.Workflow, error)
	SaveStepStatus(models.StepsStatus) (models.StepsStatus, error)
	FindStepStatusByServiceRequestID(serviceRequestID uuid.UUID) ([]models.StepsStatus, error)
	FindStepStatusByServiceRequestIDAndStatus(serviceRequestID uuid.UUID, status models.Status) ([]models.StepsStatus, error)
	FindStepStatusByServiceRequestIDAndStepIDAndStatus(serviceRequestID uuid.UUID, stepID int, status models.Status) (models.StepsStatus, error)
	FindAllStepStatusByServiceRequestIDAndStepID(serviceRequestID uuid.UUID, stepID int) ([]models.StepsStatus, error)
	GetWorkflows(pageNumber int, pageSize int, sortBy models.SortByFields) ([]models.Workflow, int, error)
	FindServiceRequestsByWorkflowName(workflowName string, pageNumber int, pageSize int) ([]models.ServiceRequest, error)
	DeleteWorkflowByName(string) error
}

var db DBInterface

func init() {
	switch config.ENV.DBDriver {
	case "postgres":
		db = &postgres{}
	}
}

// GetDB returns the initialized database implementations. Currently only postgres is implemented.
func GetDB() DBInterface {
	return db
}

// SetDb is used to update the db object with custom implementations.
// It is used in tests to override the actual db implementations with mock implementations
func SetDb(dbImpl DBInterface) {
	db = dbImpl
}
