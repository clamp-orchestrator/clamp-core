package services

import (
	"clamp-core/models"
	"fmt"
)

func SaveServiceFlow(serviceFlowReg models.Workflow) (models.Workflow, error) {
	pgServReq := serviceFlowReg.ToPGWorkflow()
	err := repo.insertQuery(&pgServReq)

	fmt.Println(pgServReq)

	if err != nil {
		panic(err)
	}
	return serviceFlowReg, err
}
