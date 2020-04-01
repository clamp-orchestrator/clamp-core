package services

import (
	"clamp-core/repository"
	"github.com/go-pg/pg/v9/orm"
	"log"
)

type Result = orm.Result

type serviceRequestRepo interface {
	selectQuery(interface{}) error
	insertQuery(interface{}) error
	whereQuery(interface{}, string, ...interface{}) error
	query(interface{}, interface{}) (Result, error)
}

type serviceRequestRepoImpl struct {
}

func (s serviceRequestRepoImpl) whereQuery(model interface{}, condition string, params ...interface{}) error {
	db := repository.GetDB()
	return db.Model(model).Where(condition, params...).Select()
}

func (s serviceRequestRepoImpl) insertQuery(model interface{}) error {
	db := repository.GetDB()
	return db.Insert(model)
}

func (s serviceRequestRepoImpl) selectQuery(model interface{}) error {
	db := repository.GetDB()
	return db.Select(model)
}

func (s serviceRequestRepoImpl) query(query interface{}, param interface{}) (Result, error) {
	db := repository.GetDB()
	result, err := db.Exec(query, param)
	log.Print("Result ", result.Model())
	log.Print("Affected ", result.RowsAffected())
	log.Print("Returned ", result.RowsReturned())
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	return result, err
}

var repo serviceRequestRepo

func init() {
	repo = serviceRequestRepoImpl{}
}
