package repository

import (
	"clamp-core/domain"
	"fmt"

	pg "github.com/go-pg/pg/v9"
	"github.com/google/uuid"
)

type pgServiceRequest struct {
	tableName    struct{} `pg:"service_requests"`
	ID           uuid.UUID
	WorkflowName string
	Status       domain.Status
}

func from(serviceReq domain.ServiceRequest) pgServiceRequest {
	return pgServiceRequest{
		ID:           serviceReq.ID,
		WorkflowName: serviceReq.WorkflowName,
		Status:       serviceReq.Status,
	}
}

func (pgServReq pgServiceRequest) to() domain.ServiceRequest {
	return domain.ServiceRequest{
		ID:           pgServReq.ID,
		WorkflowName: pgServReq.WorkflowName,
		Status:       pgServReq.Status,
	}
}

//FindByID is
func FindByID(serviceReq domain.ServiceRequest) {
	db := pg.Connect(&pg.Options{
		User:     "clamp",
		Password: "clamppass",
		Database: "clampdev",
	})
	defer db.Close()

	pgServReq := from(serviceReq)
	err := db.Select(&pgServReq)

	fmt.Println(pgServReq)

	if err != nil {
		panic(err)
	}
}

func SaveServiceRequest(serviceReq domain.ServiceRequest) domain.ServiceRequest {
	db := pg.Connect(&pg.Options{
		User:     "clamp",
		Password: "clamppass",
		Database: "clampdev",
	})
	defer db.Close()

	pgServReq := from(serviceReq)
	err := db.Insert(&pgServReq)

	fmt.Println(pgServReq)

	if err != nil {
		panic(err)
	}
	return serviceReq
}
