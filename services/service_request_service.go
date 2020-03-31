package services

import (
	"clamp-core/models"
	"clamp-core/repository"
	"fmt"
)

//FindByID is
func FindByID(serviceReq models.ServiceRequest) {
	db := repository.GetDB()

	pgServReq := serviceReq.ToPgServiceRequest()
	err := db.Select(&pgServReq)

	fmt.Println(pgServReq)

	if err != nil {
		panic(err)
	}
}

func SaveServiceRequest(serviceReq models.ServiceRequest) models.ServiceRequest {
	db := repository.GetDB()

	pgServReq := serviceReq.ToPgServiceRequest()
	err := db.Insert(&pgServReq)

	if err != nil {
		panic(err)
	}
	return serviceReq
}
