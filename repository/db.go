package repository

import (
	"clamp-core/config"
	"clamp-core/models"
	"github.com/google/uuid"
)

type dbInterface interface {
	SaveServiceRequest(models.ServiceRequest) (models.ServiceRequest, error)
	FindServiceRequestById(uuid.UUID) (models.ServiceRequest, error)
	SaveWorkflow(models.Workflow) (models.Workflow, error)
	FindWorkflowByName(string) (models.Workflow, error)
	SaveStepStatus(models.StepsStatus) (models.StepsStatus, error)
	FindStepStatusByServiceRequestId(serviceRequestId uuid.UUID) ([]models.StepsStatus, error)
	FindStepStatusByServiceRequestIdAndStatus(serviceRequestId uuid.UUID, status models.Status) ([]models.StepsStatus, error)
	FindStepStatusByServiceRequestIdAndStepIdAndStatus(serviceRequestId uuid.UUID, stepId int, status models.Status) (models.StepsStatus, error)
	FindAllStepStatusByServiceRequestIdAndStepId(serviceRequestId uuid.UUID, stepId int) ([]models.StepsStatus, error)
	GetWorkflows(pageNumber int, pageSize int) ([]models.Workflow, error)
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
