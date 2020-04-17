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
	FindStepStatusByServiceRequestIdAndStatusOrderByCreatedAtDesc(serviceRequestId uuid.UUID, status models.Status) (models.StepsStatus, error)
}

var DB dbInterface

func init() {
	switch config.Config.DBDriver {
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
