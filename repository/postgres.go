package repository

import (
	"clamp-core/config"
	"clamp-core/models"
	"context"
	"github.com/go-pg/pg/v9"
	"github.com/google/uuid"
	"log"
	"strings"
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
	db = pg.Connect(GetPostgresOptions())
	if LogSQLQueries {
		db.AddQueryHook(dbLogger{})
	}
	return db
}

func GetPostgresOptions() *pg.Options {
	connStr := config.ENV.DBConnectionStr
	connArr := strings.Split(connStr, " ")
	var host, user, password, dbName string
	for _, conn := range connArr {
		connMap := strings.Split(conn, "=")
		switch connMap[0] {
		case "host":
			host = connMap[1]
		case "user":
			user = connMap[1]
		case "password":
			password = connMap[1]
		case "dbname":
			dbName = connMap[1]
		}
	}
	return &pg.Options{
		Addr:     host,
		User:     user,
		Password: password,
		Database: dbName,
	}
}

type postgres struct {
	db *pg.DB
}

func (p *postgres) FindServiceRequestsByWorkflowName(workflowName string, pageNumber int, pageSize int) ([]models.ServiceRequest, error) {
	var pgServiceRequests []models.PGServiceRequest
	err := p.getDb().Model(&pgServiceRequests).
		Where("workflow_name = ?", workflowName).
		Limit(pageSize).
		Offset(pageNumber).
		Select()
	var workflows []models.ServiceRequest
	for _, pgServiceRequest := range pgServiceRequests {
		workflows = append(workflows, pgServiceRequest.ToServiceRequest())
	}
	return workflows, err
}

func (p *postgres) FindAllStepStatusByServiceRequestIDAndStepID(serviceRequestID uuid.UUID, stepID int) ([]models.StepsStatus, error) {
	var pgStepStatus []models.PGStepStatus
	err := p.getDb().Model(&pgStepStatus).Where("service_request_id = ? and step_id = ?", serviceRequestID, stepID).Select()
	var stepStatuses []models.StepsStatus
	for _, status := range pgStepStatus {
		stepStatuses = append(stepStatuses, status.ToStepStatus())
	}
	return stepStatuses, err
}

func (p *postgres) FindStepStatusByServiceRequestIDAndStepIDAndStatus(serviceRequestID uuid.UUID, stepID int, status models.Status) (models.StepsStatus, error) {
	var pgStepStatus models.PGStepStatus
	var stepStatuses models.StepsStatus
	err := p.getDb().Model(&pgStepStatus).Where("service_request_id = ? and step_id = ? and status = ?", serviceRequestID, stepID, status).Select()
	if err != nil {
		return stepStatuses, err
	}
	return pgStepStatus.ToStepStatus(), err
}

func (p *postgres) FindStepStatusByServiceRequestIDAndStatus(serviceRequestID uuid.UUID, status models.Status) ([]models.StepsStatus, error) {
	var pgStepStatus []models.PGStepStatus
	var stepStatuses []models.StepsStatus
	err := p.getDb().Model(&pgStepStatus).Where("service_request_id = ? and status = ?", serviceRequestID, status).Order("created_at ASC").Select()
	if err != nil {
		return stepStatuses, err
	}
	for _, status := range pgStepStatus {
		stepStatuses = append(stepStatuses, status.ToStepStatus())
	}
	return stepStatuses, err
}

func (p *postgres) FindStepStatusByServiceRequestID(serviceRequestID uuid.UUID) ([]models.StepsStatus, error) {
	var pgStepStatus []models.PGStepStatus
	err := p.getDb().Model(&pgStepStatus).Where("service_request_id = ?", serviceRequestID).Order("created_at ASC").Select()
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

func (p *postgres) FindServiceRequestByID(serviceRequestID uuid.UUID) (models.ServiceRequest, error) {
	pgServiceRequest := &models.PGServiceRequest{ID: serviceRequestID}
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

func (p *postgres) GetWorkflows(pageNumber int, pageSize int) ([]models.Workflow, error) {
	var pgWorkflows []models.PGWorkflow
	err := p.getDb().Model(&pgWorkflows).
		Limit(pageSize).
		Offset(pageNumber).
		Select()
	var workflows []models.Workflow
	for _, pgWorkflow := range pgWorkflows {
		workflows = append(workflows, pgWorkflow.ToWorkflow())
	}
	return workflows, err
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
