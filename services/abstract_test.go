package services

import (
	"clamp-core/models"
	"clamp-core/repository"

	"github.com/google/uuid"
)

type mockDBServiceRequest struct{}
type mockDBWorkFlow struct{}
type mockDBStepStatus struct{}

func (m mockDBWorkFlow) DeleteWorkflowByName(workflowName string) error {
	return deleteWorkflowByNameMock(workflowName)
}

func (m mockDBServiceRequest) FindServiceRequestsByWorkflowName(workflowName string, pageNumber int, pageSize int) ([]*models.ServiceRequest, error) {
	return findServiceRequestsByWorkflowName(workflowName, pageNumber, pageSize)
}

func (m mockDBWorkFlow) GetWorkflows(pageNumber int, pageSize int, sortBy models.SortByFields) ([]*models.Workflow, int, error) {
	return getWorkflowsMock(pageNumber, pageSize, sortBy)
}

func (m mockDBStepStatus) FindStepStatusByServiceRequestIDAndStepIDAndStatus(serviceRequestID uuid.UUID, stepID int, status models.Status) (*models.StepsStatus, error) {
	return findStepStatusByServiceRequestIDAndStepIDAndStatusMock(serviceRequestID, stepID, status)
}

func (m mockDBStepStatus) FindStepStatusByServiceRequestIDAndStepNameAndStatus(serviceRequestID uuid.UUID, stepName string, status models.Status) (*models.StepsStatus, error) {
	return findStepStatusByServiceRequestIDAndStepNameAndStatusMock(serviceRequestID, stepName, status)
}

var findServiceRequestsByWorkflowName func(workflowName string, pageNumber int, pageSize int) ([]*models.ServiceRequest, error)
var saveServiceRequestMock func(serReq *models.ServiceRequest) (*models.ServiceRequest, error)
var saveStepStatusMock func(stepStatus *models.StepsStatus) (*models.StepsStatus, error)
var SaveWorkflowMock func(workflow *models.Workflow) (*models.Workflow, error)
var findServiceRequestByIDMock func(uuid.UUID) (*models.ServiceRequest, error)
var findWorkflowByNameMock func(workflowName string) (*models.Workflow, error)
var findStepStatusByServiceRequestIDMock func(serviceRequestID uuid.UUID) ([]*models.StepsStatus, error)
var findStepStatusByServiceRequestIDAndStatusMock func(serviceRequestID uuid.UUID, status models.Status) ([]*models.StepsStatus, error)
var findAllStepStatusByServiceRequestIDAndStepIDMock func(serviceRequestID uuid.UUID, stepID int) ([]*models.StepsStatus, error)
var findStepStatusByServiceRequestIDAndStepNameAndStatusMock func(serviceRequestID uuid.UUID, stepName string, status models.Status) (*models.StepsStatus, error)
var findStepStatusByServiceRequestIDAndStepIDAndStatusMock func(serviceRequestID uuid.UUID, stepID int, status models.Status) (*models.StepsStatus, error)
var getWorkflowsMock func(pageNumber int, pageSize int, sortBy models.SortByFields) ([]*models.Workflow, int, error)
var deleteWorkflowByNameMock func(workflowName string) error

func (m mockDBServiceRequest) SaveServiceRequest(serReq *models.ServiceRequest) (*models.ServiceRequest, error) {
	return saveServiceRequestMock(serReq)
}

func (m mockDBServiceRequest) FindServiceRequestByID(id uuid.UUID) (*models.ServiceRequest, error) {
	return findServiceRequestByIDMock(id)
}

func (m mockDBWorkFlow) SaveWorkflow(workflow *models.Workflow) (*models.Workflow, error) {
	return SaveWorkflowMock(workflow)
}

func (m mockDBWorkFlow) FindWorkflowByName(workflowName string) (*models.Workflow, error) {
	return findWorkflowByNameMock(workflowName)
}

func (m mockDBStepStatus) SaveStepStatus(stepStatus *models.StepsStatus) (*models.StepsStatus, error) {
	return saveStepStatusMock(stepStatus)
}

func (m mockDBStepStatus) FindStepStatusByServiceRequestID(serviceRequestID uuid.UUID) ([]*models.StepsStatus, error) {
	return findStepStatusByServiceRequestIDMock(serviceRequestID)
}

func (m mockDBStepStatus) FindStepStatusByServiceRequestIDAndStatus(serviceRequestID uuid.UUID, status models.Status) ([]*models.StepsStatus, error) {
	return findStepStatusByServiceRequestIDAndStatusMock(serviceRequestID, status)
}

func (m mockDBStepStatus) FindAllStepStatusByServiceRequestIDAndStepID(serviceRequestID uuid.UUID, stepID int) ([]*models.StepsStatus, error) {
	return findAllStepStatusByServiceRequestIDAndStepIDMock(serviceRequestID, stepID)
}

func (m mockDBStepStatus) Ping() error {
	return nil
}

func init() {
	repository.SetServiceRequestRepository(&mockDBServiceRequest{})
	repository.SetWorkFlowRepository(&mockDBWorkFlow{})
	repository.SetStepStatusRepository(&mockDBStepStatus{})
}
