package repository

import (
	"clamp-core/models"
	"fmt"
	"sort"
	"sync"

	"github.com/go-pg/pg/v9"
	"github.com/google/uuid"
)

type memoryDB struct {
	sync.Mutex

	workflows                  map[string]*models.Workflow
	serviceRequests            map[uuid.UUID]*models.ServiceRequest
	serviceRequestStepStatuses map[uuid.UUID][]*models.StepsStatus
}

func NewMemoryDB() DBInterface {
	mdb := &memoryDB{
		workflows:                  make(map[string]*models.Workflow),
		serviceRequests:            make(map[uuid.UUID]*models.ServiceRequest),
		serviceRequestStepStatuses: make(map[uuid.UUID][]*models.StepsStatus),
	}

	return mdb
}

func (m *memoryDB) FindWorkflowByName(workflowName string) (*models.Workflow, error) {
	m.Lock()
	defer m.Unlock()

	workflow, ok := m.workflows[workflowName]
	if !ok {
		return nil, pg.ErrNoRows
	}

	return workflow, nil
}

func (m *memoryDB) GetWorkflows(pageNumber int, pageSize int, sortBy models.SortByFields) ([]*models.Workflow, int, error) {
	m.Lock()
	defer m.Unlock()

	var workflows []*models.Workflow

	for _, workflow := range m.workflows {
		workflows = append(workflows, workflow)
	}

	sort.Slice(workflows, func(i, j int) bool {
		return workflows[i].Name < workflows[j].Name
	})

	offset := pageNumber * pageSize
	return workflows[offset : offset+pageSize], len(workflows), nil
}

func (m *memoryDB) SaveWorkflow(workflow *models.Workflow) (*models.Workflow, error) {
	m.Lock()
	defer m.Unlock()

	_, exist := m.workflows[workflow.Name]
	if exist {
		return nil, fmt.Errorf("workflow with the name '%s' already exist", workflow.Name)
	}

	workflowCopy := *workflow
	m.workflows[workflow.Name] = &workflowCopy

	return m.workflows[workflow.Name], nil
}

func (m *memoryDB) DeleteWorkflowByName(workflowName string) error {
	m.Lock()
	defer m.Unlock()

	delete(m.workflows, workflowName)

	return nil
}

func (m *memoryDB) FindServiceRequestsByWorkflowName(workflowName string, pageNumber int, pageSize int) ([]*models.ServiceRequest, error) {
	m.Lock()
	defer m.Unlock()

	var serviceRequests []*models.ServiceRequest
	for _, serviceRequest := range m.serviceRequests {
		if serviceRequest.WorkflowName == workflowName {
			serviceRequests = append(serviceRequests, serviceRequest)
		}
	}

	return serviceRequests, nil
}

func (m *memoryDB) SaveServiceRequest(serReq *models.ServiceRequest) (*models.ServiceRequest, error) {
	m.Lock()
	defer m.Unlock()

	sr := *serReq
	m.serviceRequests[sr.ID] = &sr

	return &sr, nil
}

func (m *memoryDB) FindServiceRequestByID(id uuid.UUID) (*models.ServiceRequest, error) {
	m.Lock()
	defer m.Unlock()

	sr, ok := m.serviceRequests[id]
	if !ok {
		return nil, pg.ErrNoRows
	}

	return sr, nil
}

func (m *memoryDB) FindStepStatusByServiceRequestIDAndStepIDAndStatus(serviceRequestID uuid.UUID, stepID int, status models.Status) (*models.StepsStatus, error) {
	m.Lock()
	defer m.Unlock()

	stepStatuses, ok := m.serviceRequestStepStatuses[serviceRequestID]
	if !ok {
		return nil, pg.ErrNoRows
	}

	for _, stepStatus := range stepStatuses {
		if stepStatus.StepID == stepID && stepStatus.Status == status {
			return stepStatus, nil
		}
	}

	return nil, pg.ErrNoRows
}

func (m *memoryDB) FindStepStatusByServiceRequestIDAndStepNameAndStatus(serviceRequestID uuid.UUID, stepName string, status models.Status) (*models.StepsStatus, error) {
	m.Lock()
	defer m.Unlock()

	stepStatuses, ok := m.serviceRequestStepStatuses[serviceRequestID]
	if !ok {
		return nil, pg.ErrNoRows
	}

	for _, stepStatus := range stepStatuses {
		if stepStatus.StepName == stepName && stepStatus.Status == status {
			return stepStatus, nil
		}
	}

	return nil, pg.ErrNoRows
}

func (m *memoryDB) FindStepStatusByServiceRequestID(serviceRequestID uuid.UUID) ([]*models.StepsStatus, error) {
	m.Lock()
	defer m.Unlock()

	stepStatuses, ok := m.serviceRequestStepStatuses[serviceRequestID]
	if !ok {
		return nil, pg.ErrNoRows
	}

	var filteredStepStatuses []*models.StepsStatus

	for _, stepStatus := range stepStatuses {
		if stepStatus.ServiceRequestID == serviceRequestID {
			filteredStepStatuses = append(filteredStepStatuses, stepStatus)
		}
	}

	return filteredStepStatuses, nil
}

func (m *memoryDB) FindStepStatusByServiceRequestIDAndStatus(serviceRequestID uuid.UUID, status models.Status) ([]*models.StepsStatus, error) {
	m.Lock()
	defer m.Unlock()

	stepStatuses, ok := m.serviceRequestStepStatuses[serviceRequestID]
	if !ok {
		return nil, pg.ErrNoRows
	}

	var filteredStepStatuses []*models.StepsStatus

	for _, stepStatus := range stepStatuses {
		if stepStatus.ServiceRequestID == serviceRequestID && stepStatus.Status == status {
			filteredStepStatuses = append(filteredStepStatuses, stepStatus)
		}
	}

	return filteredStepStatuses, nil
}

func (m *memoryDB) FindAllStepStatusByServiceRequestIDAndStepID(serviceRequestID uuid.UUID, stepID int) ([]*models.StepsStatus, error) {
	m.Lock()
	defer m.Unlock()

	stepStatuses, ok := m.serviceRequestStepStatuses[serviceRequestID]
	if !ok {
		return nil, pg.ErrNoRows
	}

	var filteredStepStatuses []*models.StepsStatus

	for _, stepStatus := range stepStatuses {
		if stepStatus.ServiceRequestID == serviceRequestID && stepStatus.StepID == stepID {
			filteredStepStatuses = append(filteredStepStatuses, stepStatus)
		}
	}

	return filteredStepStatuses, nil
}

func (m *memoryDB) SaveStepStatus(stepStatus *models.StepsStatus) (*models.StepsStatus, error) {
	m.Lock()
	defer m.Unlock()

	_, ok := m.serviceRequestStepStatuses[stepStatus.ServiceRequestID]
	if !ok {
		m.serviceRequestStepStatuses[stepStatus.ServiceRequestID] = make([]*models.StepsStatus, 0)
	}

	ss := *stepStatus
	m.serviceRequestStepStatuses[stepStatus.ServiceRequestID] = append(m.serviceRequestStepStatuses[stepStatus.ServiceRequestID], &ss)

	return &ss, nil
}

func (m *memoryDB) Ping() error {
	return nil
}
