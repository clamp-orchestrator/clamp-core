package services

import (
	"clamp-core/models"
)

func SaveServiceFlow(serviceFlowReg models.Workflow) models.Workflow {
	//db := repository.GetDB()
	//
	//pgServReq := serviceFlowReg.ToPGWorkflow()
	//err := db.Insert(&pgServReq)

	//fmt.Println(pgServReq)
	//
	//if err != nil {
	//	panic(err)
	//}
	return serviceFlowReg
}
