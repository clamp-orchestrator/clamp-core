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

	queryMock = func(model interface{}, query interface{}, param interface{}) (result Result, err error) {
		step1Time := time.Date(2020, time.April, 07, 16, 32, 00, 00, time.UTC)
		step2Time := time.Date(2020, time.April, 07, 16, 32, 00, 20000000, time.UTC)

		test := model.(*[]models.StepsStatus)
		*test = make([]models.StepsStatus, 2)
		(*test)[0].WorkflowName = stepsStatusReq.WorkflowName
		(*test)[0].ID = stepsStatusReq.ID
		(*test)[0].Status = stepsStatusReq.Status
		(*test)[0].StepName = stepsStatusReq.StepName
		(*test)[0].TotalTimeInMs = stepsStatusReq.TotalTimeInMs
		(*test)[0].CreatedAt = step1Time
		(*test)[1].WorkflowName = stepsStatusReq.WorkflowName
		(*test)[1].ID = "2"
		(*test)[1].Status = stepsStatusReq.Status
		(*test)[1].StepName = "step2"
		(*test)[1].TotalTimeInMs = stepsStatusReq.TotalTimeInMs
		(*test)[1].CreatedAt = step2Time
		return result, err
	}
	resp, err := FindStepStatusByServiceRequestId(stepsStatusReq.ServiceRequestId)
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

	queryMock = func(model interface{}, query interface{}, param interface{}) (result Result, err error) {
		return result, errors.New("select query failed")
	}
	_, err = FindStepStatusByServiceRequestId(stepsStatusReq.ServiceRequestId)
	assert.NotNil(t, err)
}
