package repository

import (
	"clamp-core/config"
	"clamp-core/models"

	"github.com/google/uuid"
)

// DBInterface provides a collection of method signatures that needs to be implemented for a specific database.
type StepStatusRepository interface {
	SaveStepStatus(*models.StepsStatus) (*models.StepsStatus, error)
	FindStepStatusByServiceRequestID(serviceRequestID uuid.UUID) ([]*models.StepsStatus, error)
	FindStepStatusByServiceRequestIDAndStatus(serviceRequestID uuid.UUID, status models.Status) ([]*models.StepsStatus, error)
	FindStepStatusByServiceRequestIDAndStepIDAndStatus(
		serviceRequestID uuid.UUID, stepID int, status models.Status) (*models.StepsStatus, error)
	FindAllStepStatusByServiceRequestIDAndStepID(serviceRequestID uuid.UUID, stepID int) ([]*models.StepsStatus, error)
}

var stepStatusRepository StepStatusRepository

func init() {
	switch config.ENV.DBDriver {
	case "postgres":
		stepStatusRepository = &stepstatusrepositorypostgres{}
	}
}

// GetDB returns the initialized database implementations. Currently only postgres is implemented.
func GetStepStatusRepository() StepStatusRepository {
	return stepStatusRepository
}

// SetDB is used to update the db object with custom implementations.
// It is used in tests to override the actual db implementations with mock implementations
func SetStepStatusRepository(stepStatusRepository StepStatusRepository) {
	stepStatusRepository = stepStatusRepository
}
