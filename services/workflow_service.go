package services

import (
	"clamp-core/models"
	"fmt"
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
func FindWorkflowByName(workflowName string) {
	fmt.Print("Inside Find work flow function ",workflowName)
	query := "select id from workflows where name = ?"
	res, err := repo.query(query, workflowName)
	if err != nil {
		panic(err)
	}
	fmt.Print("result is ", res)
}