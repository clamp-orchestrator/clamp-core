package services

import (
	"clamp-core/models"
	"clamp-core/repository"
	"github.com/google/uuid"
)

type mockDB struct{}

var saveServiceRequestMock func(serReq models.ServiceRequest) (models.ServiceRequest, error)
var saveStepStatusMock func(stepStatus models.StepsStatus) (models.StepsStatus, error)
var SaveWorkflowMock func(workflow models.Workflow) (models.Workflow, error)
var findServiceRequestByIdMock func(uuid.UUID) (models.ServiceRequest, error)
var findWorkflowByNameMock func(workflowName string) (models.Workflow, error)
var findStepStatusByServiceRequestIdMock func(serviceRequestId uuid.UUID) ([]models.StepsStatus, error)

func (m mockDB) SaveServiceRequest(serReq models.ServiceRequest) (models.ServiceRequest, error) {
	return saveServiceRequestMock(serReq)
}

func (m mockDB) FindServiceRequestById(id uuid.UUID) (models.ServiceRequest, error) {
	return findServiceRequestByIdMock(id)
}

func (m mockDB) SaveWorkflow(workflow models.Workflow) (models.Workflow, error) {
	return SaveWorkflowMock(workflow)
}

func (m mockDB) FindWorkflowByName(workflowName string) (models.Workflow, error) {
	return findWorkflowByNameMock(workflowName)
}

func (m mockDB) SaveStepStatus(stepStatus models.StepsStatus) (models.StepsStatus, error) {
	return saveStepStatusMock(stepStatus)
}

func (m mockDB) FindStepStatusByServiceRequestId(serviceRequestId uuid.UUID) ([]models.StepsStatus, error) {
	return findStepStatusByServiceRequestIdMock(serviceRequestId)
}

func init() {
	repository.SetDb(&mockDB{})
}
