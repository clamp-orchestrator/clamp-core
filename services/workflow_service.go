package services

import (
	"clamp-core/models"
	"fmt"
	"github.com/google/uuid"
)

func SaveServiceFlow(serviceFlowReg models.Workflow) (models.Workflow, error) {
	pgServReq := serviceFlowReg.ToPGWorkflow()
	err := repo.insertQuery(&pgServReq)

	fmt.Println("",pgServReq)

	if err != nil {
		panic(err)
	}
	return serviceFlowReg, err
}

//FindServiceRequestByID is
func FindWorkflowByName(workflowName uuid.UUID) (models.Workflow, error) {
	workflowReq := models.Workflow{ID: workflowName}
	fmt.Println("Workflow request is -- ",workflowReq)
	pgWorkflowReq := workflowReq.ToPGWorkflow()
	fmt.Println("Request is -- ",pgWorkflowReq)
	err := repo.selectQuery(&pgWorkflowReq)
	if err != nil {
		panic(err)
	}
	fmt.Print( "Finally ---", pgWorkflowReq)
	return workflowReq,err
}