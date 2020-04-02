package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

//Workflow is a structure to store the service request details
type Workflow struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
	Steps       []Step `json:"steps"`
}

type Step struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Mode      string `json:"mode"`
	URL       string `json:"url"`
	Transform bool   `json:"transform"`
	Enabled   bool   `json:"enabled"`
}

//Create a new work flow for a given service flow and return service flow details
func CreateWorkflow(workflowRequest Workflow) Workflow {
	return newServiceFlow(workflowRequest)
}

func newServiceFlow(workflow Workflow) Workflow {
	return Workflow{Id: workflow.Id, Name: workflow.Name, Description: workflow.Description, Enabled: true, CreatedAt:time.Time{}, UpdatedAt:time.Time{}, Steps: workflow.Steps}
}

type PGWorkflow struct {
	tableName   struct{} `pg:"workflows"`
	Id          string
	Name        string
	Description string
	Enabled     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Steps       []Step
}

func (workflow Workflow) ToPGWorkflow() PGWorkflow {
	return PGWorkflow{
		Id:          workflow.Id,
		Name:        workflow.Name,
		Description: workflow.Description,
		Enabled:     workflow.Enabled,
		CreatedAt: workflow.CreatedAt,
		UpdatedAt: workflow.UpdatedAt,
		Steps:       workflow.Steps,
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
