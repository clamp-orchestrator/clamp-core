package services

import (
	"clamp-core/models"
	"clamp-core/repository"

	"github.com/google/uuid"
)

type mockDB struct{}

func (m mockDB) DeleteWorkflowByName(workflowName string) error {
	return deleteWorkflowByNameMock(workflowName)
}

func (m mockDB) FindServiceRequestsByWorkflowName(workflowName string, pageNumber int, pageSize int, sortBy models.SortByFields) ([]models.ServiceRequest, int, error) {
	return findServiceRequestsByWorkflowName(workflowName, pageNumber, pageSize, sortBy)
}

func (m mockDB) GetWorkflows(pageNumber int, pageSize int, sortBy models.SortByFields) ([]models.Workflow, int, error) {
	return getWorkflowsMock(pageNumber, pageSize, sortBy)
}

func (m mockDB) FindStepStatusByServiceRequestIDAndStepIDAndStatus(serviceRequestID uuid.UUID, stepID int, status models.Status) (models.StepsStatus, error) {
	return findStepStatusByServiceRequestIDAndStepIDAndStatusMock(serviceRequestID, stepID, status)
}

func (m mockDB) FindStepStatusByServiceRequestIDAndStepNameAndStatus(serviceRequestID uuid.UUID, stepName string, status models.Status) (models.StepsStatus, error) {
	return findStepStatusByServiceRequestIDAndStepNameAndStatusMock(serviceRequestID, stepName, status)
}

var findServiceRequestsByWorkflowName func(workflowName string, pageNumber int, pageSize int, sortBy models.SortByFields) ([]models.ServiceRequest, int, error)
var saveServiceRequestMock func(serReq models.ServiceRequest) (models.ServiceRequest, error)
var saveStepStatusMock func(stepStatus models.StepsStatus) (models.StepsStatus, error)
var SaveWorkflowMock func(workflow models.Workflow) (models.Workflow, error)
var findServiceRequestByIDMock func(uuid.UUID) (models.ServiceRequest, error)
var findWorkflowByNameMock func(workflowName string) (models.Workflow, error)
var findStepStatusByServiceRequestIDMock func(serviceRequestID uuid.UUID) ([]models.StepsStatus, error)
var findStepStatusByServiceRequestIDAndStatusMock func(serviceRequestID uuid.UUID, status models.Status) ([]models.StepsStatus, error)
var findAllStepStatusByServiceRequestIDAndStepIDMock func(serviceRequestID uuid.UUID, stepID int) ([]models.StepsStatus, error)
var findStepStatusByServiceRequestIDAndStepNameAndStatusMock func(serviceRequestID uuid.UUID, stepName string, status models.Status) (models.StepsStatus, error)
var findStepStatusByServiceRequestIDAndStepIDAndStatusMock func(serviceRequestID uuid.UUID, stepID int, status models.Status) (models.StepsStatus, error)
var getWorkflowsMock func(pageNumber int, pageSize int, sortBy models.SortByFields) ([]models.Workflow, int, error)
var deleteWorkflowByNameMock func(workflowName string) error

func (m mockDB) SaveServiceRequest(serReq models.ServiceRequest) (models.ServiceRequest, error) {
	return saveServiceRequestMock(serReq)
}

func (m mockDB) FindServiceRequestByID(id uuid.UUID) (models.ServiceRequest, error) {
	return findServiceRequestByIDMock(id)
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

func (m mockDB) FindStepStatusByServiceRequestID(serviceRequestID uuid.UUID) ([]models.StepsStatus, error) {
	return findStepStatusByServiceRequestIDMock(serviceRequestID)
}

func (m mockDB) FindStepStatusByServiceRequestIDAndStatus(serviceRequestID uuid.UUID, status models.Status) ([]models.StepsStatus, error) {
	return findStepStatusByServiceRequestIDAndStatusMock(serviceRequestID, status)
}

func (m mockDB) FindAllStepStatusByServiceRequestIDAndStepID(serviceRequestID uuid.UUID, stepID int) ([]models.StepsStatus, error) {
	return findAllStepStatusByServiceRequestIDAndStepIDMock(serviceRequestID, stepID)
}

func (m mockDB) Ping() error {
	return nil
}

func init() {
	repository.SetDb(&mockDB{})
}
