package repository

import (
	"clamp-core/models"
	"context"
	"github.com/go-pg/pg/v9"
	"github.com/google/uuid"
	"log"
	"sync"
)

const LogSQLQueries bool = true

var singletonOnce sync.Once

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	query, err := q.FormattedQuery()
	log.Printf("[PSQL] Query: %v, Error: %v", query, err)
	return nil
}

func connectDB() (db *pg.DB) {
	db = pg.Connect(&pg.Options{
		User:     "clamp",
		Password: "clamppass",
		Database: "clampdev",
	})
	if LogSQLQueries {
		db.AddQueryHook(dbLogger{})
	}
	return db
}

type postgres struct {
	db *pg.DB
}

func (p *postgres) FindStepStatusByServiceRequestIdAndStepNameAndStatus(serviceRequestId uuid.UUID, stepName string, status models.Status) (models.StepsStatus, error) {
	var pgStepStatus models.PGStepStatus
	var stepStatuses models.StepsStatus
	err := p.getDb().Model(&pgStepStatus).Where("service_request_id = ? and step_name = ? and status = ?", serviceRequestId, stepName, status).Select()
	if err != nil {
		return stepStatuses, err
	}
	return pgStepStatus.ToStepStatus(), err
}

func (p *postgres) FindStepStatusByServiceRequestIdAndStatusOrderByCreatedAtDesc(serviceRequestId uuid.UUID, status models.Status) (models.StepsStatus, error) {
	var pgStepStatus []models.PGStepStatus
	var stepStatuses models.StepsStatus
	err := p.getDb().Model(&pgStepStatus).Where("service_request_id = ? and status = ?", serviceRequestId, status).Order("created_at DESC").Select()
	if err != nil {
		return stepStatuses, err
	}
	return pgStepStatus[0].ToStepStatus(), err
}

func (p *postgres) FindStepStatusByServiceRequestId(serviceRequestId uuid.UUID) ([]models.StepsStatus, error) {
	var pgStepStatus []models.PGStepStatus
	err := p.getDb().Model(&pgStepStatus).Where("service_request_id = ?", serviceRequestId).Select()
	var stepStatuses []models.StepsStatus
	for _, status := range pgStepStatus {
		stepStatuses = append(stepStatuses, status.ToStepStatus())
	}
	return stepStatuses, err
}

func (p *postgres) SaveStepStatus(stepStatus models.StepsStatus) (models.StepsStatus, error) {
	pgStepStatusReq := stepStatus.ToPgStepStatus()
	err := p.getDb().Insert(&pgStepStatusReq)
	return pgStepStatusReq.ToStepStatus(), err
}

func (p *postgres) FindWorkflowByName(workflowName string) (models.Workflow, error) {
	pgWorkflow := new(models.PGWorkflow)
	err := p.getDb().Model(pgWorkflow).Where("name = ?", workflowName).Select()
	return (*pgWorkflow).ToWorkflow(), err
}

func (p *postgres) SaveWorkflow(workflowReq models.Workflow) (models.Workflow, error) {
	pgWorkflow := workflowReq.ToPGWorkflow()
	log.Printf("pgworfklow: %v", pgWorkflow)
	err := p.getDb().Insert(&pgWorkflow)
	return pgWorkflow.ToWorkflow(), err
}

func (p *postgres) FindServiceRequestById(serviceRequestId uuid.UUID) (models.ServiceRequest, error) {
	pgServiceRequest := &models.PGServiceRequest{ID: serviceRequestId}
	err := p.getDb().Select(pgServiceRequest)
	if err != nil {
		panic(err)
	}
	return (*pgServiceRequest).ToServiceRequest(), err
}

func (p *postgres) SaveServiceRequest(serviceReq models.ServiceRequest) (models.ServiceRequest, error) {
	pgServReq := serviceReq.ToPgServiceRequest()
	db := p.getDb()
	err := db.Insert(&pgServReq)
	return pgServReq.ToServiceRequest(), err
}

func (p *postgres) getDb() *pg.DB {
	singletonOnce.Do(func() {
		log.Println("Connecting to DB")
		p.db = connectDB()
	})
	return p.db
}

func (p *postgres) closeDB() {
	if p.db != nil {
		log.Println("Disconnecting from DB")
		p.db.Close()
	}
}
