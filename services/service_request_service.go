package services

import (
	"clamp-core/models"
	"clamp-core/repository"
	"fmt"
	"github.com/google/uuid"
)

type pgServiceRequest struct {
	tableName    struct{} `pg:"service_requests"`
	ID           uuid.UUID
	WorkflowName string
	Status       models.Status
}

func from(serviceReq models.ServiceRequest) pgServiceRequest {
	return pgServiceRequest{
		ID:           serviceReq.ID,
		WorkflowName: serviceReq.WorkflowName,
		Status:       serviceReq.Status,
	}
}

func (pgServReq pgServiceRequest) to() models.ServiceRequest {
	return models.ServiceRequest{
		ID:           pgServReq.ID,
		WorkflowName: pgServReq.WorkflowName,
		Status:       pgServReq.Status,
	}
}

//FindByID is
func FindByID(serviceReq models.ServiceRequest) {
	db := repository.GetDB()

	pgServReq := from(serviceReq)
	err := db.Select(&pgServReq)

	fmt.Println(pgServReq)

	if err != nil {
		panic(err)
	}
}

func SaveServiceRequest(serviceReq models.ServiceRequest) models.ServiceRequest {
	db := repository.GetDB()

	pgServReq := from(serviceReq)
	err := db.Insert(&pgServReq)

	if err != nil {
		panic(err)
	}
	return serviceReq
}
