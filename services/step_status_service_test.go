package services

import (
	"clamp-core/models"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func prepareStepsStatus() models.StepsStatus {
	stepsStatus := models.StepsStatus{
		ID:               "1",
		ServiceRequestId: uuid.New(),
		WorkflowName:     workflowName,
		Status:           models.STATUS_COMPLETED,
		CreatedAt:        time.Now(),
		StepName:         "Testing",
		TotalTimeInMs:    10,
	}
	return stepsStatus
}
func TestSaveStepsStatus(t *testing.T) {

	stepsStatusReq := prepareStepsStatus()
	saveStepStatusMock = func(stepStatus models.StepsStatus) (status models.StepsStatus, err error) {
		return stepStatus, nil
	}
	response, err := SaveStepStatus(stepsStatusReq)
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.NotNil(t, response.ID)
	assert.Equal(t, stepsStatusReq.StepName, response.StepName, fmt.Sprintf("Expected Step name to be %s but was %s", stepsStatusReq.StepName, response.StepName))
	assert.Equal(t, stepsStatusReq.TotalTimeInMs, response.TotalTimeInMs, fmt.Sprintf("Expected Total time in ms to be %d but was %d", stepsStatusReq.TotalTimeInMs, response.TotalTimeInMs))
	assert.Equal(t, stepsStatusReq.Status, response.Status, fmt.Sprintf("Expected Step status to be %s but was %s", stepsStatusReq.Status, response.Status))

	saveStepStatusMock = func(stepStatus models.StepsStatus) (status models.StepsStatus, err error) {
		return status, errors.New("insertion failed")
	}
	response, err = SaveStepStatus(stepsStatusReq)
	assert.NotNil(t, err)
}

func TestFindStepStatusByServiceRequestId(t *testing.T) {
	stepsStatusReq := prepareStepsStatus()
	findStepStatusByServiceRequestIdMock = func(serviceRequestId uuid.UUID) (statuses []models.StepsStatus, err error) {
		step1Time := time.Date(2020, time.April, 07, 16, 32, 00, 00, time.UTC)
		step2Time := time.Date(2020, time.April, 07, 16, 32, 00, 20000000, time.UTC)

		statuses = make([]models.StepsStatus, 4)
		statuses[0].WorkflowName = stepsStatusReq.WorkflowName
		statuses[0].ID = stepsStatusReq.ID
		statuses[0].Status = models.STATUS_STARTED
		statuses[0].StepName = stepsStatusReq.StepName
		statuses[0].TotalTimeInMs = stepsStatusReq.TotalTimeInMs
		statuses[0].CreatedAt = step1Time
		statuses[1].WorkflowName = stepsStatusReq.WorkflowName
		statuses[1].ID = stepsStatusReq.ID
		statuses[1].Status = stepsStatusReq.Status
		statuses[1].StepName = stepsStatusReq.StepName
		statuses[1].TotalTimeInMs = stepsStatusReq.TotalTimeInMs
		statuses[1].CreatedAt = step1Time
		statuses[2].WorkflowName = stepsStatusReq.WorkflowName
		statuses[2].ID = "2"
		statuses[2].Status = models.STATUS_STARTED
		statuses[2].StepName = "step2"
		statuses[2].TotalTimeInMs = stepsStatusReq.TotalTimeInMs
		statuses[2].CreatedAt = step2Time
		statuses[3].WorkflowName = stepsStatusReq.WorkflowName
		statuses[3].ID = "2"
		statuses[3].Status = stepsStatusReq.Status
		statuses[3].StepName = "step2"
		statuses[3].TotalTimeInMs = stepsStatusReq.TotalTimeInMs
		statuses[3].CreatedAt = step2Time
		return statuses, err
	}

	stepsStatus, err := FindStepStatusByServiceRequestId(stepsStatusReq.ServiceRequestId)
	workflow := models.Workflow{
		Name:        stepsStatusReq.WorkflowName,
		Description: "",
		Enabled:     false,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Steps:       make([]models.Step, 2),
	}
	resp := PrepareStepStatusResponse(stepsStatusReq.ServiceRequestId, workflow, stepsStatus)
	assert.Nil(t, err)
	assert.Equal(t, stepsStatusReq.WorkflowName, resp.WorkflowName)
	assert.Equal(t, models.STATUS_COMPLETED, resp.Status)
	//assert.Equal(t, stepsStatusReq.ServiceRequestId, resp.ServiceRequestId)
	assert.Equal(t, int64(20), resp.TotalTimeInMs)
	assert.NotNil(t, resp.ServiceRequestId)
	assert.NotNil(t, resp.Steps)
	assert.Equal(t, models.STATUS_COMPLETED, resp.Steps[1].Status)
	assert.Equal(t, stepsStatusReq.StepName, resp.Steps[0].Name)
	assert.Equal(t, stepsStatusReq.TotalTimeInMs, resp.Steps[0].TimeTaken)
	assert.Equal(t, models.STATUS_COMPLETED, resp.Steps[3].Status)
	assert.Equal(t, "step2", resp.Steps[2].Name)
	assert.Equal(t, stepsStatusReq.TotalTimeInMs, resp.Steps[2].TimeTaken)

	findStepStatusByServiceRequestIdMock = func(serviceRequestId uuid.UUID) (statuses []models.StepsStatus, err error) {
		return statuses, errors.New("select query failed")
	}
	_, err = FindStepStatusByServiceRequestId(stepsStatusReq.ServiceRequestId)
	assert.NotNil(t, err)
}

func TestShouldReturnStatusCompletedForAllStepsCompleted(t *testing.T) {
	findStepStatusByServiceRequestIdMock = func(serviceRequestId uuid.UUID) (statuses []models.StepsStatus, err error) {
		step1Time := time.Date(2020, time.April, 07, 16, 32, 00, 00, time.UTC)

		statuses = make([]models.StepsStatus, 4)
		workflowName := "testWF"
		statuses[0].WorkflowName = workflowName
		statuses[0].ID = "1"
		statuses[0].Status = models.STATUS_STARTED
		statuses[0].StepName = "step1"
		statuses[0].TotalTimeInMs = 10
		statuses[0].CreatedAt = step1Time
		statuses[1].WorkflowName = workflowName
		statuses[1].ID = "2"
		statuses[1].Status = models.STATUS_COMPLETED
		statuses[1].StepName = "step1"
		statuses[1].TotalTimeInMs = 20
		statuses[1].CreatedAt = step1Time.Add(time.Second * 10)

		statuses[2].WorkflowName = workflowName
		statuses[2].ID = "3"
		statuses[2].Status = models.STATUS_STARTED
		statuses[2].StepName = "step2"
		statuses[2].TotalTimeInMs = 10
		statuses[2].CreatedAt = step1Time.Add(time.Second * 20)
		statuses[3].WorkflowName = workflowName
		statuses[3].ID = "4"
		statuses[3].Status = models.STATUS_COMPLETED
		statuses[3].StepName = "step2"
		statuses[3].TotalTimeInMs = 20
		statuses[3].CreatedAt = step1Time.Add(time.Second * 30)
		return statuses, err
	}

	serviceReqID := uuid.New()
	stepsStatus, err := FindStepStatusByServiceRequestId(serviceReqID)
	workflow := models.Workflow{
		Name:        workflowName,
		Description: "",
		Enabled:     false,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Steps:       make([]models.Step, 2),
	}
	resp := PrepareStepStatusResponse(serviceReqID, workflow, stepsStatus)
	assert.Nil(t, err)
	assert.Equal(t, models.STATUS_COMPLETED, resp.Status)
}

func TestShouldReturnStatusFailed(t *testing.T) {
	findStepStatusByServiceRequestIdMock = func(serviceRequestId uuid.UUID) (statuses []models.StepsStatus, err error) {
		step1Time := time.Date(2020, time.April, 07, 16, 32, 00, 00, time.UTC)

		statuses = make([]models.StepsStatus, 4)
		statuses[0].WorkflowName = "testWF"
		statuses[0].ID = "1"
		statuses[0].Status = models.STATUS_STARTED
		statuses[0].StepName = "step1"
		statuses[0].TotalTimeInMs = 10
		statuses[0].CreatedAt = step1Time
		statuses[1].WorkflowName = "testWF"
		statuses[1].ID = "2"
		statuses[1].Status = models.STATUS_COMPLETED
		statuses[1].StepName = "step1"
		statuses[1].TotalTimeInMs = 20
		statuses[1].CreatedAt = step1Time.Add(time.Second * 10)

		statuses[2].WorkflowName = "testWF"
		statuses[2].ID = "3"
		statuses[2].Status = models.STATUS_STARTED
		statuses[2].StepName = "step2"
		statuses[2].TotalTimeInMs = 10
		statuses[2].CreatedAt = step1Time.Add(time.Second * 20)
		statuses[3].WorkflowName = "testWF"
		statuses[3].ID = "4"
		statuses[3].Status = models.STATUS_FAILED
		statuses[3].StepName = "step2"
		statuses[3].TotalTimeInMs = 20
		statuses[3].CreatedAt = step1Time.Add(time.Second * 30)
		return statuses, err
	}
	serviceReqID := uuid.New()
	stepsStatus, err := FindStepStatusByServiceRequestId(serviceReqID)
	workflow := models.Workflow{
		Name:        workflowName,
		Description: "",
		Enabled:     false,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Steps:       make([]models.Step, 2),
	}
	resp := PrepareStepStatusResponse(serviceReqID, workflow, stepsStatus)
	assert.Nil(t, err)
	assert.Equal(t, models.STATUS_FAILED, resp.Status)
}

func TestShouldReturnStatusInprogress(t *testing.T) {
	findStepStatusByServiceRequestIdMock = func(serviceRequestId uuid.UUID) (statuses []models.StepsStatus, err error) {
		step1Time := time.Date(2020, time.April, 07, 16, 32, 00, 00, time.UTC)

		statuses = make([]models.StepsStatus, 4)
		statuses[0].WorkflowName = "testWF"
		statuses[0].ID = "1"
		statuses[0].Status = models.STATUS_STARTED
		statuses[0].StepName = "step1"
		statuses[0].TotalTimeInMs = 10
		statuses[0].CreatedAt = step1Time
		statuses[1].WorkflowName = "testWF"
		statuses[1].ID = "2"
		statuses[1].Status = models.STATUS_COMPLETED
		statuses[1].StepName = "step1"
		statuses[1].TotalTimeInMs = 20
		statuses[1].CreatedAt = step1Time.Add(time.Second * 10)

		statuses[2].WorkflowName = "testWF"
		statuses[2].ID = "3"
		statuses[2].Status = models.STATUS_STARTED
		statuses[2].StepName = "step2"
		statuses[2].TotalTimeInMs = 10
		statuses[2].CreatedAt = step1Time.Add(time.Second * 20)
		return statuses, err
	}

	serviceReqID := uuid.New()
	stepsStatus, err := FindStepStatusByServiceRequestId(serviceReqID)
	workflow := models.Workflow{
		Name:        workflowName,
		Description: "",
		Enabled:     false,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Steps:       make([]models.Step, 2),
	}
	resp := PrepareStepStatusResponse(serviceReqID, workflow, stepsStatus)
	assert.Nil(t, err)
	assert.Equal(t, models.STATUS_INPROGRESS, resp.Status)
}

func TestFindStepStatusByServiceRequestIdAndStatusOrderByCreatedAtDesc(t *testing.T) {
	stepsStatusReq := prepareStepsStatus()
	findStepStatusByServiceRequestIdAndStatusOrderByCreatedAtDescMock = func(serviceRequestId uuid.UUID, status models.Status) (statuses models.StepsStatus, err error) {
		step1Time := time.Date(2020, time.April, 07, 16, 32, 00, 00, time.UTC)

		statuses = models.StepsStatus{}
		statuses.WorkflowName = "testWF"
		statuses.ID = "1"
		statuses.Status = models.STATUS_STARTED
		statuses.StepName = "step1"
		statuses.TotalTimeInMs = 10
		statuses.CreatedAt = step1Time
		return statuses, err
	}

	stepsStatus, err := FindStepStatusByServiceRequestIdAndStatusOrderByCreatedAtDesc(stepsStatusReq.ServiceRequestId, models.STATUS_STARTED)

	assert.Nil(t, err)
	assert.Equal(t, stepsStatusReq.WorkflowName, stepsStatus.WorkflowName)
	assert.Equal(t, models.STATUS_STARTED, stepsStatus.Status)
	assert.Equal(t, int64(10), stepsStatus.TotalTimeInMs)
	assert.NotNil(t, stepsStatus.ServiceRequestId)

	assert.Equal(t, "step1", stepsStatus.StepName)
	assert.Equal(t, stepsStatusReq.TotalTimeInMs, stepsStatus.TotalTimeInMs)

	findStepStatusByServiceRequestIdAndStatusOrderByCreatedAtDescMock = func(serviceRequestId uuid.UUID, status models.Status) (statuses models.StepsStatus, err error) {
		return statuses, errors.New("select query failed")
	}
	_, err = FindStepStatusByServiceRequestIdAndStatusOrderByCreatedAtDesc(stepsStatusReq.ServiceRequestId, models.STATUS_STARTED)
	assert.NotNil(t, err)
}

func TestFindStepStatusByServiceRequestIdAndStepIdAndStatus(t *testing.T) {
	stepsStatusReq := prepareStepsStatus()
	findStepStatusByServiceRequestIdAndStatusOrderByCreatedAtDescMock = func(serviceRequestId uuid.UUID, status models.Status) (statuses models.StepsStatus, err error) {
		step1Time := time.Date(2020, time.April, 07, 16, 32, 00, 00, time.UTC)

		statuses = models.StepsStatus{}
		statuses.WorkflowName = "testWF"
		statuses.ID = "1"
		statuses.Status = models.STATUS_STARTED
		statuses.StepName = "step1"
		statuses.StepId = 1
		statuses.TotalTimeInMs = 10
		statuses.CreatedAt = step1Time
		return statuses, err
	}

	stepsStatus, err := FindStepStatusByServiceRequestIdAndStatusOrderByCreatedAtDesc(stepsStatusReq.ServiceRequestId, models.STATUS_STARTED)

	assert.Nil(t, err)
	assert.Equal(t, stepsStatusReq.WorkflowName, stepsStatus.WorkflowName)
	assert.Equal(t, models.STATUS_STARTED, stepsStatus.Status)
	assert.Equal(t, int64(10), stepsStatus.TotalTimeInMs)
	assert.NotNil(t, stepsStatus.ServiceRequestId)

	assert.Equal(t, "step1", stepsStatus.StepName)
	assert.Equal(t, stepsStatusReq.TotalTimeInMs, stepsStatus.TotalTimeInMs)

	findStepStatusByServiceRequestIdAndStatusOrderByCreatedAtDescMock = func(serviceRequestId uuid.UUID, status models.Status) (statuses models.StepsStatus, err error) {
		return statuses, errors.New("select query failed")
	}
	_, err = FindStepStatusByServiceRequestIdAndStatusOrderByCreatedAtDesc(stepsStatusReq.ServiceRequestId, models.STATUS_STARTED)
	assert.NotNil(t, err)
}

func TestShouldReturnStatusCompletedIfOneStepSkipped(t *testing.T) {
	findStepStatusByServiceRequestIdMock = func(serviceRequestId uuid.UUID) (statuses []models.StepsStatus, err error) {
		step1Time := time.Date(2020, time.April, 07, 16, 32, 00, 00, time.UTC)

		statuses = make([]models.StepsStatus, 4)
		statuses[0].WorkflowName = "testWF"
		statuses[0].ID = "1"
		statuses[0].Status = models.STATUS_STARTED
		statuses[0].StepName = "step1"
		statuses[0].TotalTimeInMs = 10
		statuses[0].CreatedAt = step1Time
		statuses[1].WorkflowName = "testWF"
		statuses[1].ID = "2"
		statuses[1].Status = models.STATUS_SKIPPED
		statuses[1].StepName = "step1"
		statuses[1].TotalTimeInMs = 20
		statuses[1].CreatedAt = step1Time.Add(time.Second * 10)

		statuses[2].WorkflowName = "testWF"
		statuses[2].ID = "3"
		statuses[2].Status = models.STATUS_STARTED
		statuses[2].StepName = "step2"
		statuses[2].TotalTimeInMs = 10
		statuses[2].CreatedAt = step1Time.Add(time.Second * 20)
		statuses[3].WorkflowName = "testWF"
		statuses[3].ID = "4"
		statuses[3].Status = models.STATUS_COMPLETED
		statuses[3].StepName = "step2"
		statuses[3].TotalTimeInMs = 20
		statuses[3].CreatedAt = step1Time.Add(time.Second * 30)
		return statuses, err
	}
	serviceReqID := uuid.New()
	stepsStatus, err := FindStepStatusByServiceRequestId(serviceReqID)
	workflow := models.Workflow{
		Name:        workflowName,
		Description: "",
		Enabled:     false,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Steps:       make([]models.Step, 2),
	}
	resp := PrepareStepStatusResponse(serviceReqID, workflow, stepsStatus)
	assert.Nil(t, err)
	assert.Equal(t, models.STATUS_COMPLETED, resp.Status)
}
