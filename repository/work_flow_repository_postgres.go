package repository

import (
	"clamp-core/models"
	"errors"

	"github.com/prometheus/common/log"
)

type workflowrepositorypostgres struct {
}

var keyReferences = map[string]string{"id": "id", "created_at": "created_at", "name": "name"}

func (workflowrepository *workflowrepositorypostgres) FindWorkflowByName(workflowName string) (*models.Workflow, error) {
	pgWorkflow := new(models.PGWorkflow)

	err := pgDB.GetDB().Model(pgWorkflow).Where("name = ?", workflowName).Select()
	return pgWorkflow.ToWorkflow(), err
}

func (workflowrepository *workflowrepositorypostgres) DeleteWorkflowByName(workflowName string) error {
	_, err := pgDB.GetDB().Model((*models.PGWorkflow)(nil)).Where("name = ?", workflowName).Delete()
	return err
}

func (workflowrepository *workflowrepositorypostgres) SaveWorkflow(workflowReq *models.Workflow) (*models.Workflow, error) {
	pgWorkflow := workflowReq.ToPGWorkflow()
	log.Debugf("pgworfklow: %v", pgWorkflow)
	err := pgDB.GetDB().Insert(pgWorkflow)
	return pgWorkflow.ToWorkflow(), err
}

func (workflowrepository *workflowrepositorypostgres) GetWorkflows(pageNumber int, pageSize int, sortFields models.SortByFields) ([]*models.Workflow, int, error) {
	var pgWorkflows []models.PGWorkflow
	query := pgDB.GetDB().Model(&pgWorkflows)
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
