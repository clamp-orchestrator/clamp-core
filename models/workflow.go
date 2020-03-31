package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
)

//Workflow is a structure to store the service request details
type Workflow struct {
	ID          uuid.UUID `json:"id"`
	ServiceFlow ServiceFlow
}
type ServiceFlow struct {
	Description string `json:"description"`
	FlowMode    string `json:"flowMode"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Steps       Steps  `json:"steps"`
}

type Steps struct {
	Step []Step `json:"step"`
}

type Step struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

//Create a new work flow for a given service flow and return service flow details
func CreateWorkflow(serviceFlowRequest Workflow) Workflow {
	return newServiceFlow(serviceFlowRequest)
}

func newServiceFlow(workflow Workflow) Workflow {
	return Workflow{ID: uuid.New(), ServiceFlow: workflow.ServiceFlow}
}

type PGWorkflow struct {
	tableName   struct{} `pg:"workflows"`
	ServiceFlow ServiceFlow
	ID           uuid.UUID
}

func (serviceFlow Workflow) ToPGWorkflow() PGWorkflow {
	return PGWorkflow{
		ID:           serviceFlow.ID,
		ServiceFlow: serviceFlow.ServiceFlow,
	}
}

// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (workflow PGWorkflow) Value() (driver.Value, error) {
	return json.Marshal(workflow)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (workflow *PGWorkflow) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &workflow)
}
