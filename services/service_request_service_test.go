package services

import (
	"clamp-core/models"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var insertQueryMock func(model interface{}) error
var selectQueryMock func(model interface{}) error

type mockServiceRequestRepoImpl struct {
}

func (s mockServiceRequestRepoImpl) insertQuery(model interface{}) error {
	return insertQueryMock(model)
}

func (s mockServiceRequestRepoImpl) selectQuery(model interface{}) error {
	return selectQueryMock(model)
}

func TestSaveServiceRequest(t *testing.T) {
	serviceReq := models.ServiceRequest{
		ID:           uuid.UUID{},
		WorkflowName: "TESTING",
		Status:       models.STATUS_NEW,
	}
	repo = mockServiceRequestRepoImpl{}

	insertQueryMock = func(model interface{}) error {
		return nil
	}
	request, err := SaveServiceRequest(serviceReq)
	assert.NotNil(t, request)
	assert.Nil(t, err)

	insertQueryMock = func(model interface{}) error {
		return errors.New("insertion failed")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SaveServiceRequest should have panicked!")
		}
	}()
	request, err = SaveServiceRequest(serviceReq)
}

func TestFindByID(t *testing.T) {
	serviceReq := models.ServiceRequest{
		ID: uuid.UUID{},
	}
	repo = mockServiceRequestRepoImpl{}

	selectQueryMock = func(model interface{}) error {
		serviceReq.WorkflowName = "TEST_WF"
		serviceReq.Status = models.STATUS_COMPLETED
		return nil
	}
	FindServiceRequestByID(&serviceReq)
	assert.Equal(t, "TEST_WF", serviceReq.WorkflowName)
	assert.Equal(t, models.STATUS_COMPLETED, serviceReq.Status)

	selectQueryMock = func(model interface{}) error {
		return errors.New("select query failed")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("FindServiceRequestByID should have panicked!")
		}
	}()
	FindServiceRequestByID(&serviceReq)
}
