package services

import (
	"clamp-core/models"
	"clamp-core/repository"
)

type serviceRequestRepo interface {
	selectQuery(interface{}) error
	insertQuery(interface{}) error
}

type serviceRequestRepoImpl struct {
}

func (s serviceRequestRepoImpl) insertQuery(model interface{}) error {
	db := repository.GetDB()
	return db.Insert(model)
}

func (s serviceRequestRepoImpl) selectQuery(model interface{}) error {
	db := repository.GetDB()
	return db.Select(model)
}

var repo serviceRequestRepo

func init() {
	repo = serviceRequestRepoImpl{}
}

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
