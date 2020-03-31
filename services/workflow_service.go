package services

import (
	"clamp-core/models"
	"clamp-core/repository"
	"fmt"
)

func SaveServiceFlow(serviceFlowReg models.Workflow) models.Workflow {
	db := repository.GetDB()

	pgServReq := serviceFlowReg.ToPGWorkflow()
	err := db.Insert(&pgServReq)

	fmt.Println(pgServReq)

	if err != nil {
		panic(err)
	}
	return serviceFlowReg
}
