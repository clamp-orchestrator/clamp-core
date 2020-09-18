package services

import (
	"clamp-core/models"
	"clamp-core/repository"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSaveServiceRequest(t *testing.T) {
	serviceReq := models.ServiceRequest{
		ID:           uuid.UUID{},
		WorkflowName: "TESTING",
		Status:       models.STATUS_NEW,
	}

	saveServiceRequestMock = func(serReq models.ServiceRequest) (request models.ServiceRequest, err error) {
		return serReq, nil
	}
	request, err := SaveServiceRequest(serviceReq)
	assert.NotNil(t, request)
	assert.Nil(t, err)

	saveServiceRequestMock = func(serReq models.ServiceRequest) (request models.ServiceRequest, err error) {
		return serReq, errors.New("insertion failed")
	}
	serviceReq.WorkflowName = ""
	request, err = SaveServiceRequest(serviceReq)
	assert.NotNil(t, err)
	assert.Equal(t, "insertion failed", err.Error())
}

func TestShouldFailToSaveServiceRequestAndThrowError(t *testing.T) {
	serviceReq := models.ServiceRequest{
		ID:           uuid.UUID{},
		WorkflowName: "TESTING",
		Status:       models.STATUS_NEW,
	}

	saveServiceRequestMock = func(serReq models.ServiceRequest) (request models.ServiceRequest, err error) {
		return models.ServiceRequest{}, errors.New("insertion failed")
	}
	serviceReq.WorkflowName = ""
	request, err := SaveServiceRequest(serviceReq)
	assert.Equal(t, models.ServiceRequest{}, request)
	assert.NotNil(t, err)
	assert.Equal(t, "insertion failed", err.Error())
}

func TestFindByID(t *testing.T) {
	repository.SetDb(&mockDB{})
	serviceReq := models.ServiceRequest{
		ID: uuid.UUID{},
	}
	findServiceRequestByIDMock = func(id uuid.UUID) (request models.ServiceRequest, err error) {
		request.WorkflowName = "TEST_WF"
		request.Status = models.STATUS_COMPLETED
		return request, nil
	}
	resp, err := FindServiceRequestByID(serviceReq.ID)
	assert.Nil(t, err)
	assert.Equal(t, "TEST_WF", resp.WorkflowName)
	assert.Equal(t, models.STATUS_COMPLETED, resp.Status)

	findServiceRequestByIDMock = func(id uuid.UUID) (request models.ServiceRequest, err error) {
		return request, errors.New("select query failed")
	}
	_, err = FindServiceRequestByID(serviceReq.ID)
	assert.NotNil(t, err)
	assert.Equal(t, "select query failed", err.Error())
}

func TestFindServiceRequestsByWorkflowName(t *testing.T) {
	serviceReq := models.ServiceRequest{
		ID: uuid.UUID{},
	}
	findServiceRequestsByWorkflowName = func(workflowName string, pageNumber int, pageSize int) ([]models.ServiceRequest, error) {
		return []models.ServiceRequest{serviceReq}, nil
	}
	resp, err := findServiceRequestsByWorkflowName("test", 1, 1)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(resp))
}
