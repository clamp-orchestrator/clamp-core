package repository

import (
	"clamp-core/config"
	"clamp-core/models"
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/go-pg/pg/v9"
	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
)

// reference human readable keys to DB key values
var keyReferences = map[string]string{"id": "id", "created_at": "created_at", "name": "name"}

var singletonOnce sync.Once

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	query, err := q.FormattedQuery()
	log.Infof("[PSQL] Query: %v, Error: %v", query, err)
	return nil
}

func connectDB() (db *pg.DB) {
	db = pg.Connect(GetPostgresOptions())
	if config.ENV.EnableSQLQueriesLog {
		db.AddQueryHook(dbLogger{})
	}
	return db
}

// GetPostgresOptions returns connection details for postgres DB
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

func (p *postgres) FindServiceRequestsByWorkflowName(workflowName string, pageNumber int, pageSize int) ([]*models.ServiceRequest, error) {
	var pgServiceRequests []models.PGServiceRequest
	err := p.getDB().Model(&pgServiceRequests).
		Where("workflow_name = ?", workflowName).
		Offset(pageSize * pageNumber).
		Limit(pageSize).
		Select()
	var workflows []*models.ServiceRequest
	if err == nil {
		for _, pgServiceRequest := range pgServiceRequests {
			workflows = append(workflows, pgServiceRequest.ToServiceRequest())
		}
	}
	return workflows, err
}

func (p *postgres) FindAllStepStatusByServiceRequestIDAndStepID(serviceRequestID uuid.UUID, stepID int) ([]*models.StepsStatus, error) {
	var pgStepStatus []models.PGStepStatus
	err := p.getDB().Model(&pgStepStatus).Where("service_request_id = ? and step_id = ?", serviceRequestID, stepID).Select()
	var stepStatuses []*models.StepsStatus
	if err == nil {
		for i := range pgStepStatus {
			stepStatuses = append(stepStatuses, pgStepStatus[i].ToStepStatus())
		}
	}
	return stepStatuses, err
}

func (p *postgres) FindStepStatusByServiceRequestIDAndStepIDAndStatus(
	serviceRequestID uuid.UUID, stepID int, status models.Status) (*models.StepsStatus, error) {
	var pgStepStatus models.PGStepStatus
	var stepStatuses models.StepsStatus
	err := p.getDB().Model(&pgStepStatus).Where("service_request_id = ? and step_id = ? and status = ?",
		serviceRequestID, stepID, status).Select()
	if err != nil {
		return &stepStatuses, err
	}
	return pgStepStatus.ToStepStatus(), err
}

func (p *postgres) FindStepStatusByServiceRequestIDAndStatus(
	serviceRequestID uuid.UUID, status models.Status) ([]*models.StepsStatus, error) {
	var pgStepStatus []models.PGStepStatus
	var stepStatuses []*models.StepsStatus
	err := p.getDB().Model(&pgStepStatus).Where("service_request_id = ? and status = ?", serviceRequestID, status).
		Order("created_at ASC").Select()
	if err != nil {
		return stepStatuses, err
	}
	for i := range pgStepStatus {
		stepStatuses = append(stepStatuses, pgStepStatus[i].ToStepStatus())
	}
	return stepStatuses, err
}

func (p *postgres) FindStepStatusByServiceRequestID(serviceRequestID uuid.UUID) ([]*models.StepsStatus, error) {
	var pgStepStatus []models.PGStepStatus
	err := p.getDB().Model(&pgStepStatus).Where("service_request_id = ?", serviceRequestID).Order("created_at ASC").Select()
	var stepStatuses []*models.StepsStatus
	if err == nil {
		for i := range pgStepStatus {
			stepStatuses = append(stepStatuses, pgStepStatus[i].ToStepStatus())
		}
	}
	return stepStatuses, err
}

func (p *postgres) SaveStepStatus(stepStatus *models.StepsStatus) (*models.StepsStatus, error) {
	pgStepStatusReq := stepStatus.ToPgStepStatus()
	err := p.getDB().Insert(pgStepStatusReq)
	return pgStepStatusReq.ToStepStatus(), err
}

func (p *postgres) FindWorkflowByName(workflowName string) (*models.Workflow, error) {
	pgWorkflow := new(models.PGWorkflow)
	err := p.getDB().Model(pgWorkflow).Where("name = ?", workflowName).Select()
	return pgWorkflow.ToWorkflow(), err
}

func (p *postgres) DeleteWorkflowByName(workflowName string) error {
	_, err := p.getDB().Model((*models.PGWorkflow)(nil)).Where("name = ?", workflowName).Delete()
	return err
}

func (p *postgres) SaveWorkflow(workflowReq *models.Workflow) (*models.Workflow, error) {
	pgWorkflow := workflowReq.ToPGWorkflow()
	log.Debugf("pgworfklow: %v", pgWorkflow)
	err := p.getDB().Insert(pgWorkflow)
	return pgWorkflow.ToWorkflow(), err
}

func (p *postgres) FindServiceRequestByID(serviceRequestID uuid.UUID) (*models.ServiceRequest, error) {
	pgServiceRequest := &models.PGServiceRequest{ID: serviceRequestID}
	err := p.getDB().Select(pgServiceRequest)
	if err != nil {
		panic(err)
	}
	return pgServiceRequest.ToServiceRequest(), err
}

func (p *postgres) SaveServiceRequest(serviceReq *models.ServiceRequest) (*models.ServiceRequest, error) {
	pgServReq := serviceReq.ToPgServiceRequest()
	db := p.getDB()
	err := db.Insert(pgServReq)
	return pgServReq.ToServiceRequest(), err
}

func (p *postgres) GetWorkflows(pageNumber int, pageSize int, sortFields models.SortByFields) ([]*models.Workflow, int, error) {
	var pgWorkflows []models.PGWorkflow
	query := p.getDB().Model(&pgWorkflows)
	for _, sortField := range sortFields {
		reference, found := keyReferences[sortField.Key]
		if !found {
			return []*models.Workflow{}, 0, errors.New("undefined key reference used")
		}
		order := sortField.Order
		if found {
			query = query.Order(reference + " " + order)
		}
	}
	totalWorkflowsCount, err := query.Offset(pageSize * (pageNumber - 1)).
		Limit(pageSize).SelectAndCount()
	if err != nil {
		return []*models.Workflow{}, 0, err
	}
	var workflows []*models.Workflow
	for i := range pgWorkflows {
		workflows = append(workflows, pgWorkflows[i].ToWorkflow())
	}
	return workflows, totalWorkflowsCount, err
}

func (p *postgres) Ping() error {
	_, err := p.getDB().Exec("SELECT 1")
	return err
}

func (p *postgres) getDB() *pg.DB {
	singletonOnce.Do(func() {
		log.Info("Connecting to DB")
		p.db = connectDB()
	})
	return p.db
}

func (p *postgres) closeDB() {
	if p.db != nil {
		log.Info("Disconnecting from DB")
		p.db.Close()
	}
}
