package repository

import (
	"clamp-core/models"

	"github.com/google/uuid"
)

type MockDB struct {
	FindServiceRequestsByWorkflowNameFunc                        func(workflowName string, pageNumber int, pageSize int) ([]*models.ServiceRequest, error)
	SaveServiceRequestMockFunc                                   func(serReq *models.ServiceRequest) (*models.ServiceRequest, error)
	SaveStepStatusMockFunc                                       func(stepStatus *models.StepsStatus) (*models.StepsStatus, error)
	SaveWorkflowMockFunc                                         func(workflow *models.Workflow) (*models.Workflow, error)
	FindServiceRequestByIDMockFunc                               func(uuid.UUID) (*models.ServiceRequest, error)
	FindWorkflowByNameMockFunc                                   func(workflowName string) (*models.Workflow, error)
	FindStepStatusByServiceRequestIDMockFunc                     func(serviceRequestID uuid.UUID) ([]*models.StepsStatus, error)
	FindStepStatusByServiceRequestIDAndStatusMockFunc            func(serviceRequestID uuid.UUID, status models.Status) ([]*models.StepsStatus, error)
	FindAllStepStatusByServiceRequestIDAndStepIDMockFunc         func(serviceRequestID uuid.UUID, stepID int) ([]*models.StepsStatus, error)
	FindStepStatusByServiceRequestIDAndStepNameAndStatusMockFunc func(serviceRequestID uuid.UUID, stepName string, status models.Status) (*models.StepsStatus, error)
	FindStepStatusByServiceRequestIDAndStepIDAndStatusMockFunc   func(serviceRequestID uuid.UUID, stepID int, status models.Status) (*models.StepsStatus, error)
	GetWorkflowsMockFunc                                         func(pageNumber int, pageSize int, sortBy models.SortByFields) ([]*models.Workflow, int, error)
	DeleteWorkflowByNameMockFunc                                 func(workflowName string) error
}

func (m *MockDB) DeleteWorkflowByName(workflowName string) error {
	return m.DeleteWorkflowByNameMockFunc(workflowName)
}

func (m *MockDB) FindServiceRequestsByWorkflowName(workflowName string, pageNumber int, pageSize int) ([]*models.ServiceRequest, error) {
	return m.FindServiceRequestsByWorkflowNameFunc(workflowName, pageNumber, pageSize)
}

func (m *MockDB) GetWorkflows(pageNumber int, pageSize int, sortBy models.SortByFields) ([]*models.Workflow, int, error) {
	return m.GetWorkflowsMockFunc(pageNumber, pageSize, sortBy)
}

func (m *MockDB) FindStepStatusByServiceRequestIDAndStepIDAndStatus(serviceRequestID uuid.UUID, stepID int, status models.Status) (*models.StepsStatus, error) {
	return m.FindStepStatusByServiceRequestIDAndStepIDAndStatusMockFunc(serviceRequestID, stepID, status)
}

func (m *MockDB) FindStepStatusByServiceRequestIDAndStepNameAndStatus(serviceRequestID uuid.UUID, stepName string, status models.Status) (*models.StepsStatus, error) {
	return m.FindStepStatusByServiceRequestIDAndStepNameAndStatusMockFunc(serviceRequestID, stepName, status)
}

func (m *MockDB) SaveServiceRequest(serReq *models.ServiceRequest) (*models.ServiceRequest, error) {
	return m.SaveServiceRequestMockFunc(serReq)
}

func (m *MockDB) FindServiceRequestByID(id uuid.UUID) (*models.ServiceRequest, error) {
	return m.FindServiceRequestByIDMockFunc(id)
}

func (m *MockDB) SaveWorkflow(workflow *models.Workflow) (*models.Workflow, error) {
	return m.SaveWorkflowMockFunc(workflow)
}

func (m *MockDB) FindWorkflowByName(workflowName string) (*models.Workflow, error) {
	return m.FindWorkflowByNameMockFunc(workflowName)
}

func (m *MockDB) SaveStepStatus(stepStatus *models.StepsStatus) (*models.StepsStatus, error) {
	return m.SaveStepStatusMockFunc(stepStatus)
}

func (m *MockDB) FindStepStatusByServiceRequestID(serviceRequestID uuid.UUID) ([]*models.StepsStatus, error) {
	return m.FindStepStatusByServiceRequestIDMockFunc(serviceRequestID)
}

func (m *MockDB) FindStepStatusByServiceRequestIDAndStatus(serviceRequestID uuid.UUID, status models.Status) ([]*models.StepsStatus, error) {
	return m.FindStepStatusByServiceRequestIDAndStatusMockFunc(serviceRequestID, status)
}

func (m *MockDB) FindAllStepStatusByServiceRequestIDAndStepID(serviceRequestID uuid.UUID, stepID int) ([]*models.StepsStatus, error) {
	return m.FindAllStepStatusByServiceRequestIDAndStepIDMockFunc(serviceRequestID, stepID)
}

func (m *MockDB) Ping() error {
	return nil
}
