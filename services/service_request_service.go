package services

import (
	"clamp-core/models"
)

//FindServiceRequestByID is
func FindServiceRequestByID(serviceReq *models.ServiceRequest) {
	pgServReq := serviceReq.ToPgServiceRequest()
	err := repo.selectQuery(&pgServReq)
	if err != nil {
		panic(err)
	}
}

func SaveServiceRequest(serviceReq models.ServiceRequest) (models.ServiceRequest, error) {
	pgServReq := serviceReq.ToPgServiceRequest()
	err := repo.insertQuery(&pgServReq)

	if err != nil {
		panic(err)
	}
	return serviceReq, err
}
