package repository

import (
	"clamp-core/config"
	"clamp-core/models"

	"github.com/google/uuid"
)

type dbInterface interface {
	SaveServiceRequest(models.ServiceRequest) (models.ServiceRequest, error)
	FindServiceRequestByID(uuid.UUID) (models.ServiceRequest, error)
	SaveWorkflow(models.Workflow) (models.Workflow, error)
	FindWorkflowByName(string) (models.Workflow, error)
	SaveStepStatus(models.StepsStatus) (models.StepsStatus, error)
	FindStepStatusByServiceRequestID(serviceRequestID uuid.UUID) ([]models.StepsStatus, error)
	FindStepStatusByServiceRequestIDAndStatus(serviceRequestID uuid.UUID, status models.Status) ([]models.StepsStatus, error)
	FindStepStatusByServiceRequestIDAndStepIDAndStatus(serviceRequestID uuid.UUID, stepID int, status models.Status) (models.StepsStatus, error)
	FindAllStepStatusByServiceRequestIDAndStepID(serviceRequestID uuid.UUID, stepID int) ([]models.StepsStatus, error)
	GetWorkflows(pageNumber int, pageSize int, sortBy models.SortByFields) ([]models.Workflow, error)
	FindServiceRequestsByWorkflowName(workflowName string, pageNumber int, pageSize int) ([]models.ServiceRequest, error)
}

var DB dbInterface

func init() {
	switch config.ENV.DBDriver {
	case "postgres":
		DB = &postgres{}
	}
}

func GetDB() dbInterface {
	return DB
}

func SetDb(dbImpl dbInterface) {
	DB = dbImpl
}
