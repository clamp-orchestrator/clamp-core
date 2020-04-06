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
var queryMock func(model interface{}, query interface{}, param interface{}) (Result, error)
var whereQueryMock func(model interface{}, cond string, params ...interface{}) error

type mockGenericRepoImpl struct {
}

func (s mockGenericRepoImpl) whereQuery(models interface{}, cond string, params ...interface{}) error {
	return whereQueryMock(models, cond, params)
}

func (s mockGenericRepoImpl) query(model interface{}, query interface{}, params interface{}) (Result, error) {
	return queryMock(model, query, params)
}

func (s mockGenericRepoImpl) insertQuery(model interface{}) error {
	return insertQueryMock(model)
}

func (s mockGenericRepoImpl) selectQuery(model interface{}) error {
	return selectQueryMock(model)
}

func TestSaveServiceRequest(t *testing.T) {
	serviceReq := models.ServiceRequest{
		ID:           uuid.UUID{},
		WorkflowName: "TESTING",
		Status:       models.STATUS_NEW,
	}
	repo = mockGenericRepoImpl{}

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
	repo = mockGenericRepoImpl{}

	whereQueryMock = func(model interface{}, cond string, params ...interface{}) error {
		serviceReq := model.(*models.ServiceRequest)
		serviceReq.WorkflowName = "TEST_WF"
		serviceReq.Status = models.STATUS_COMPLETED
		return nil
	}
	resp, err := FindServiceRequestByID(serviceReq.ID)
	assert.Nil(t, err)
	assert.Equal(t, "TEST_WF", resp.WorkflowName)
	assert.Equal(t, models.STATUS_COMPLETED, resp.Status)

	whereQueryMock = func(model interface{}, cond string, params ...interface{}) error {
		return errors.New("select query failed")
	}
	_, err = FindServiceRequestByID(serviceReq.ID)
	assert.NotNil(t, err)
}
