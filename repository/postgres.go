package repository

import (
	"clamp-core/domain"
	"database/sql/driver"
	"encoding/json"
	"errors"
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

type pgWorkflow struct {
	tableName    struct{} `pg:"workflows"`
	ServiceFlow domain.ServiceFlow
}

func from(serviceReq domain.ServiceRequest) pgServiceRequest {
	return pgServiceRequest{
		ID:           serviceReq.ID,
		WorkflowName: serviceReq.WorkflowName,
		Status:       serviceReq.Status,
	}
}

func fromServiceFlow(serviceFlow domain.Workflow) pgWorkflow {
	return pgWorkflow{
		//ID:           serviceFlow.ID,
		ServiceFlow: serviceFlow.ServiceFlow,
	}
}

// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (workflow pgWorkflow) Value() (driver.Value, error) {
	return json.Marshal(workflow)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (workflow *pgWorkflow) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &workflow)
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

func SaveServiceFlow(serviceFlowReg domain.Workflow) domain.Workflow {
	db := pg.Connect(&pg.Options{
		User:     "clamp",
		Password: "clamppass",
		Database: "clampdev",
	})
	defer db.Close()

	pgServReq := fromServiceFlow(serviceFlowReg)
	err := db.Insert(&pgServReq)

	fmt.Println(pgServReq)

	if err != nil {
		panic(err)
	}
	return serviceFlowReg
}
