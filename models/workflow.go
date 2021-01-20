package models

import (
	"clamp-core/config"
	"clamp-core/executors"
	"clamp-core/utils"
	"time"
)

type stepEnricher func(Step, int) Step
type stepType string

var stepEnrichmentMap map[stepType]stepEnricher = map[stepType]stepEnricher{

	"AMQP": func(step Step, stepId int) Step {
		step.ID = stepId
		step.Val.(*executors.AMQPVal).ReplyTo = config.ENV.QueueName
		step.Type = utils.AsyncStepType
		return step
	},

	"HTTP": func(step Step, stepId int) Step {
		step.ID = stepId
		if step.Type == "" {
			step.Type = utils.SyncStepType
		}
		return step
	},

	"KAFKA": func(step Step, stepId int) Step {
		step.ID = stepId
		step.Val.(*executors.KafkaVal).ReplyTo = config.ENV.KafkaConsumerTopicName
		step.Type = utils.AsyncStepType
		return step
	},
}

//Workflow is a structure to store the workflow details. A workflow record represents a definition of a workflow with multiple steps in it.
//A workflow when trigerred creates a service request. Each step in a workflow is triggered sequentially and status for workflow is tracked
//to completion (success/failure)
type Workflow struct {
	ID          string    `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Steps       []Step    `json:"steps" binding:"required,gt=0,dive"`
}

//PGWorkflow is a struct which represents the structure of workflow record in the postgres data store. The Workflow struct
//is converted to PGWorkflow while saving and vice versa during retrieval.
type PGWorkflow struct {
	tableName   struct{} `pg:"workflows"`
	ID          string
	Name        string
	Description string
	Enabled     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Steps       []Step
}

//ToPGWorkflow converts a Workflow object into a Postgres specific workflow stucture that is used for persisting to postgres specifically
func (workflow Workflow) ToPGWorkflow() PGWorkflow {
	return PGWorkflow{
		ID:          workflow.ID,
		Name:        workflow.Name,
		Description: workflow.Description,
		Enabled:     workflow.Enabled,
		CreatedAt:   workflow.CreatedAt,
		UpdatedAt:   workflow.UpdatedAt,
		Steps:       workflow.Steps,
	}
}

//ToWorkflow converts a PGWorkflow struct to a Workflow structure that is used in the application to pass around a workflow defintion.
func (pgWorkflow PGWorkflow) ToWorkflow() Workflow {
	return Workflow{
		ID:          pgWorkflow.ID,
		Name:        pgWorkflow.Name,
		Description: pgWorkflow.Description,
		Enabled:     pgWorkflow.Enabled,
		CreatedAt:   pgWorkflow.CreatedAt,
		UpdatedAt:   pgWorkflow.UpdatedAt,
		Steps:       pgWorkflow.Steps,
	}
}

//CreateWorkflow creates a new work flow with the given Workflow defintion and step details. The workflow is persisted and can be triggered.
//Once a workflow has been created it is not possible to change/edit it currently. During creation each step is assigned an id, which is unique to
//the workflow, a type and a reply to queue name for Kafka and AMQP to ensure that for async channels a reply channel is set during workflow
//creation.
func CreateWorkflow(workflowRequest Workflow) Workflow {
	stepCount := 0
	for i := 0; i < len(workflowRequest.Steps); i++ {
		stepCount++
		stepTypeName := stepType(workflowRequest.Steps[i].Mode)
		workflowRequest.Steps[i] = stepEnrichmentMap[stepTypeName](workflowRequest.Steps[i], stepCount)
		updateStepCounterForEachOfSubSteps(workflowRequest, i, stepCount)
	}
	return newServiceFlow(workflowRequest)
}

func updateStepCounterForEachOfSubSteps(workflowRequest Workflow, i int, stepCount int) {
	if workflowRequest.Steps[i].OnFailure != nil {
		stepCount = updateSubStepsIds(workflowRequest, i, stepCount)
	}
}

func updateSubStepsIds(workflowRequest Workflow, i int, stepCount int) int {
	subSteps := workflowRequest.Steps[i].OnFailure
	for subStepID := range subSteps {
		stepCount++
		subSteps[subStepID].ID = stepCount
	}
	return stepCount
}

func newServiceFlow(workflow Workflow) Workflow {
	return Workflow{
		ID:          workflow.ID,
		Name:        workflow.Name,
		Description: workflow.Description,
		Enabled:     true,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Steps:       workflow.Steps,
	}
}
