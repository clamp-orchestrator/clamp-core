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

//FindServiceRequestByID is
func FindWorkflowByName(workflowName string) (*models.Workflow, error) {
	workflowReq := models.Workflow{Name: workflowName}
	fmt.Println("Workflow request is -- ", workflowReq)
	pgWorkflowReq := workflowReq.ToPGWorkflow()
	fmt.Println("Request is -- ", pgWorkflowReq)
	//query := "select id, name from workflows where name = ?"
	//res, err := repo.query(query, workflowName)
	workflow := new(models.Workflow)
	err := repo.whereQuery(workflow, "workflow.name = ?", workflowName)
	if err != nil {
		panic(err)
	}
	fmt.Print("Finally ---", pgWorkflowReq)
	return workflow, err
}
