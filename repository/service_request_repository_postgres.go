package repository

import (
	"clamp-core/models"

	"github.com/google/uuid"
	"github.com/prometheus/common/log"
)

type servicerequestrepositorypostgres struct {
}

func (servicerequestrepository *servicerequestrepositorypostgres) SaveServiceRequest(serviceReq *models.ServiceRequest) (*models.ServiceRequest, error) {
	pgServReq := serviceReq.ToPgServiceRequest()
	db := pgDB.GetDB()
	err := db.Insert(pgServReq)
	return pgServReq.ToServiceRequest(), err
}

func (servicerequestrepository *servicerequestrepositorypostgres) FindServiceRequestsByWorkflowName(workflowName string, pageNumber int, pageSize int) ([]*models.ServiceRequest, error) {
	var pgServiceRequests []models.PGServiceRequest
	err := pgDB.GetDB().Model(&pgServiceRequests).
		Where("WORKFLOW_NAME = ?", workflowName).
		Offset(pageSize * pageNumber).
		Limit(pageSize).
		Select()
	var workflows []*models.ServiceRequest
	log.Info("workflow name %s %v", workflowName, pgServiceRequests)

	if err == nil {
		for _, pgServiceRequest := range pgServiceRequests {
			workflows = append(workflows, pgServiceRequest.ToServiceRequest())
		}
	}
	return workflows, err
}

func (servicerequestrepository *servicerequestrepositorypostgres) FindServiceRequestByID(serviceRequestID uuid.UUID) (*models.ServiceRequest, error) {
	pgServiceRequest := &models.PGServiceRequest{ID: serviceRequestID}
	err := pgDB.GetDB().Select(pgServiceRequest)
	if err != nil {
		panic(err)
	}
	return pgServiceRequest.ToServiceRequest(), err
}
