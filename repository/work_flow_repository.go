package repository

import (
	"clamp-core/config"
	"clamp-core/models"
)

// DBInterface provides a collection of method signatures that needs to be implemented for a specific database.
type WorkFlowRepository interface {
	SaveWorkflow(*models.Workflow) (*models.Workflow, error)
	FindWorkflowByName(string) (*models.Workflow, error)
	GetWorkflows(pageNumber int, pageSize int, sortBy models.SortByFields) ([]*models.Workflow, int, error)
	DeleteWorkflowByName(string) error
}

var workFlowRepository WorkFlowRepository

func init() {
	switch config.ENV.DBDriver {
	case "postgres":
		workFlowRepository = &workflowrepositorypostgres{}
	}
}

// GetDB returns the initialized database implementations. Currently only postgres is implemented.
func GetWorkFlowRepository() WorkFlowRepository {
	return workFlowRepository
}

// SetDB is used to update the db object with custom implementations.
// It is used in tests to override the actual db implementations with mock implementations
func SetWorkFlowRepository(workFlowRepositoryImpl WorkFlowRepository) {
	workFlowRepository = workFlowRepositoryImpl
}
