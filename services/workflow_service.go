package services

import (
	"clamp-core/models"
	"log"
)

func SaveWorkflow(workflowReq models.Workflow) (models.Workflow, error) {
	pgWorkflow := workflowReq.ToPGWorkflow()
	err := repo.insertQuery(&pgWorkflow)

	if err != nil {
		log.Printf("Failed to save workflow: %v, error: %s\n", pgWorkflow, err.Error())
	}
	log.Printf("Saved worflow %v", pgWorkflow)
	return pgWorkflow.ToWorkflow(), err
}

func FindWorkflowByName(workflowName string) (*models.Workflow, error) {
	workflow := new(models.Workflow)
	err := repo.whereQuery(workflow, "workflow.name = ?", workflowName)
	if err != nil {
		log.Printf("No record found with given workflow name %s, error: %s\n", workflowName, err.Error())
	}
	return workflow, err
}
