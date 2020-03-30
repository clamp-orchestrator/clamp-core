package services

import (
	"clamp-core/models"
	"clamp-core/repository"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type pgWorkflow struct {
	tableName   struct{} `pg:"workflows"`
	ServiceFlow models.ServiceFlow
}

func fromServiceFlow(serviceFlow models.Workflow) pgWorkflow {
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

func SaveServiceFlow(serviceFlowReg models.Workflow) models.Workflow {
	db := repository.GetDB()

	pgServReq := fromServiceFlow(serviceFlowReg)
	err := db.Insert(&pgServReq)

	fmt.Println(pgServReq)

	if err != nil {
		panic(err)
	}
	return serviceFlowReg
}
