package repository

import (
	"clamp-core/models"

	"github.com/google/uuid"
)

type stepstatusrepositorypostgres struct {
}

func (stepstatusrepository *stepstatusrepositorypostgres) FindAllStepStatusByServiceRequestIDAndStepID(serviceRequestID uuid.UUID, stepID int) ([]*models.StepsStatus, error) {
	var pgStepStatus []models.PGStepStatus
	err := pgDB.GetDB().Model(&pgStepStatus).Where("service_request_id = ? and step_id = ?", serviceRequestID, stepID).Select()
	var stepStatuses []*models.StepsStatus
	if err == nil {
		for i := range pgStepStatus {
			stepStatuses = append(stepStatuses, pgStepStatus[i].ToStepStatus())
		}
	}
	return stepStatuses, err
}

func (stepstatusrepository *stepstatusrepositorypostgres) FindStepStatusByServiceRequestIDAndStepIDAndStatus(
	serviceRequestID uuid.UUID, stepID int, status models.Status) (*models.StepsStatus, error) {
	var pgStepStatus models.PGStepStatus
	var stepStatuses models.StepsStatus
	err := pgDB.GetDB().Model(&pgStepStatus).Where("service_request_id = ? and step_id = ? and status = ?",
		serviceRequestID, stepID, status).Select()
	if err != nil {
		return &stepStatuses, err
	}
	return pgStepStatus.ToStepStatus(), err
}

func (stepstatusrepository *stepstatusrepositorypostgres) FindStepStatusByServiceRequestIDAndStatus(
	serviceRequestID uuid.UUID, status models.Status) ([]*models.StepsStatus, error) {
	var pgStepStatus []models.PGStepStatus
	var stepStatuses []*models.StepsStatus
	err := pgDB.GetDB().Model(&pgStepStatus).Where("service_request_id = ? and status = ?", serviceRequestID, status).
		Order("created_at ASC").Select()
	if err != nil {
		return stepStatuses, err
	}
	for i := range pgStepStatus {
		stepStatuses = append(stepStatuses, pgStepStatus[i].ToStepStatus())
	}
	return stepStatuses, err
}

func (stepstatusrepository *stepstatusrepositorypostgres) FindStepStatusByServiceRequestID(serviceRequestID uuid.UUID) ([]*models.StepsStatus, error) {
	var pgStepStatus []models.PGStepStatus
	err := pgDB.GetDB().Model(&pgStepStatus).Where("service_request_id = ?", serviceRequestID).Order("created_at ASC").Select()
	var stepStatuses []*models.StepsStatus
	if err == nil {
		for i := range pgStepStatus {
			stepStatuses = append(stepStatuses, pgStepStatus[i].ToStepStatus())
		}
	}
	return stepStatuses, err
}

func (stepstatusrepository *stepstatusrepositorypostgres) SaveStepStatus(stepStatus *models.StepsStatus) (*models.StepsStatus, error) {
	pgStepStatusReq := stepStatus.ToPgStepStatus()
	err := pgDB.GetDB().Insert(pgStepStatusReq)
	return pgStepStatusReq.ToStepStatus(), err
}
