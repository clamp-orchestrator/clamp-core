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
		Status:           models.STATUS_STARTED,
		CreatedAt:        time.Now(),
		StepName:         "Testing",
		TotalTimeInMs:    10,
	}
	return stepsStatus
}
func TestSaveStepsStatus(t *testing.T) {

	stepsStatusReq := prepareStepsStatus()
	repo = mockGenericRepoImpl{}

	insertQueryMock = func(model interface{}) error {
		return nil
	}
	response, err := SaveStepStatus(stepsStatusReq)
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.NotNil(t, response.ID)
	assert.Equal(t, stepsStatusReq.StepName, response.StepName, fmt.Sprintf("Expected Step name to be %s but was %s", stepsStatusReq.StepName, response.StepName))
	assert.Equal(t, stepsStatusReq.TotalTimeInMs, response.TotalTimeInMs, fmt.Sprintf("Expected Total time in ms to be %d but was %d", stepsStatusReq.TotalTimeInMs, response.TotalTimeInMs))
	assert.Equal(t, stepsStatusReq.Status, response.Status, fmt.Sprintf("Expected Step status to be %s but was %s", stepsStatusReq.Status, response.Status))

	insertQueryMock = func(model interface{}) error {
		return errors.New("insertion failed")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Saving Steps Status should have panicked!")
		}
	}()
	response, err = SaveStepStatus(stepsStatusReq)
}

func TestFindStepStatusByServiceRequestId(t *testing.T) {
	stepsStatusReq := prepareStepsStatus()

	repo = mockGenericRepoImpl{}

	queryMock = func(model interface{}, query interface{}, params ...interface{}) (Result, error) {
		return nil, nil
	}
	resp, err := FindStepStatusByServiceRequestId(stepsStatusReq.ServiceRequestId)
	assert.Nil(t, err)
	assert.Equal(t, stepsStatusReq.StepName, resp.Steps[0].Name)
	assert.NotNil(t, resp.Steps)

	queryMock = func(model interface{}, cond string, params ...interface{}) error {
		return errors.New("select query failed")
	}
	_, err = FindStepStatusByServiceRequestId(stepsStatusReq.ServiceRequestId)
	assert.NotNil(t, err)
}
