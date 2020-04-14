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

		statuses = make([]models.StepsStatus, 2)
		statuses[0].WorkflowName = stepsStatusReq.WorkflowName
		statuses[0].ID = stepsStatusReq.ID
		statuses[0].Status = stepsStatusReq.Status
		statuses[0].StepName = stepsStatusReq.StepName
		statuses[0].TotalTimeInMs = stepsStatusReq.TotalTimeInMs
		statuses[0].CreatedAt = step1Time
		statuses[1].WorkflowName = stepsStatusReq.WorkflowName
		statuses[1].ID = "2"
		statuses[1].Status = stepsStatusReq.Status
		statuses[1].StepName = "step2"
		statuses[1].TotalTimeInMs = stepsStatusReq.TotalTimeInMs
		statuses[1].CreatedAt = step2Time
		return statuses, err
	}

	stepsStatus, err := FindStepStatusByServiceRequestId(stepsStatusReq.ServiceRequestId)
	resp := PrepareStepStatusResponse(stepsStatus)
	assert.Nil(t, err)
	assert.Equal(t, stepsStatusReq.WorkflowName, resp.WorkflowName)
	assert.Equal(t, stepsStatusReq.Status, resp.Status)
	//assert.Equal(t, stepsStatusReq.ServiceRequestId, resp.ServiceRequestId)
	assert.Equal(t, int64(20), resp.TotalTimeInMs)
	assert.NotNil(t, resp.ServiceRequestId)
	assert.NotNil(t, resp.Steps)
	assert.Equal(t, stepsStatusReq.Status, resp.Steps[0].Status)
	assert.Equal(t, stepsStatusReq.StepName, resp.Steps[0].Name)
	assert.Equal(t, stepsStatusReq.TotalTimeInMs, resp.Steps[0].TimeTaken)
	assert.Equal(t, stepsStatusReq.Status, resp.Steps[1].Status)
	assert.Equal(t, "step2", resp.Steps[1].Name)
	assert.Equal(t, stepsStatusReq.TotalTimeInMs, resp.Steps[1].TimeTaken)

	findStepStatusByServiceRequestIdMock = func(serviceRequestId uuid.UUID) (statuses []models.StepsStatus, err error) {
		return statuses, errors.New("select query failed")
	}
	_, err = FindStepStatusByServiceRequestId(stepsStatusReq.ServiceRequestId)
	assert.NotNil(t, err)
}
