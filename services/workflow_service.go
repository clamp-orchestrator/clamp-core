package services

import (
	"clamp-core/models"
	"fmt"
)

func SaveWorkflow(workflowReq models.Workflow) (models.Workflow, error) {
	fmt.Println("Inside save service flow ", workflowReq)
	pgWorkflow := workflowReq.ToPGWorkflow()
	fmt.Println("After converting to pgworkflow request ", workflowReq)
	err := repo.insertQuery(&pgWorkflow)

	fmt.Println("Response", pgWorkflow)

	if err != nil {
		panic(err)
	}
	return workflowReq, err
}

func FindWorkflowByName(workflowName string) (*models.Workflow, error) {
	workflow := new(models.Workflow)
	err := repo.whereQuery(workflow, "workflow.name = ?", workflowName)
	if err != nil {
		fmt.Errorf("No record found with given workflow name %s", workflowName)
	}
	return workflow, err
}
