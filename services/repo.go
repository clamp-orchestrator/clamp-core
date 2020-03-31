package services

import "clamp-core/repository"

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
