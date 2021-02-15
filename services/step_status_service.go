package services

import (
	"clamp-core/models"
	"clamp-core/repository"
	"clamp-core/utils"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

var (
	serviceRequestStepNameTimeExecutorCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "service_request_with_step_name_time_counter",
		Help: "The total time taken by a step to execute",
	}, []string{"service_request_id", "step_name", "step_status"})
	serviceRequestStepNameCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "service_request_with_step_name_counter",
		Help: "Steps counter ",
	}, []string{"service_request_id", "step_name", "step_status"})
	serviceRequestStepNameTimeExecutorHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "service_request_with_step_name_histogram",
		Help: "total time taken histogram",
	})
)

func SaveStepStatus(stepStatusReq models.StepsStatus) (*models.StepsStatus, error) {
	log.Debugf("Saving step status : %v", stepStatusReq)
	stepStatusReq, err := repository.GetDB().SaveStepStatus(stepStatusReq)
	if err != nil {
		log.Errorf("Failed saving step status : %v, %s", stepStatusReq, err.Error())
	}

	serviceRequestStepNameTimeExecutorCounter.WithLabelValues(stepStatusReq.ServiceRequestID.String(),
		stepStatusReq.StepName, string(stepStatusReq.Status)).Add(float64(stepStatusReq.TotalTimeInMs))

	serviceRequestStepNameCounter.WithLabelValues(stepStatusReq.ServiceRequestID.String(),
		stepStatusReq.StepName, string(stepStatusReq.Status)).Inc()

	serviceRequestStepNameTimeExecutorHistogram.Observe(float64(stepStatusReq.TotalTimeInMs))
	return stepStatusReq, err
}

func FindStepStatusByServiceRequestID(serviceRequestID uuid.UUID) ([]*models.StepsStatus, error) {
	log.Debugf("Find step statues by request id : %s ", serviceRequestID)
	stepsStatuses, err := repository.GetDB().FindStepStatusByServiceRequestID(serviceRequestID)
	if err != nil {
		log.Errorf("No record found with given service request id %s", serviceRequestID)
		return []models.StepsStatus{}, err
	}
	return stepsStatuses, err
}

func FindStepStatusByServiceRequestIDAndStatus(serviceRequestID uuid.UUID, status models.Status) ([]*models.StepsStatus, error) {
	log.Debugf("Find step statues by request id : %s ", serviceRequestID)
	stepsStatuses, err := repository.GetDB().FindStepStatusByServiceRequestIDAndStatus(serviceRequestID, status)
	if err != nil {
		log.Errorf("No record found with given service request id %s", serviceRequestID)
		return []models.StepsStatus{}, err
	}
	return stepsStatuses, err
}

func FindAllStepStatusByServiceRequestIDAndStepID(serviceRequestID uuid.UUID, stepID int) ([]*models.StepsStatus, error) {
	log.Debugf("Find all step statues by request id : %s and step id : %d", serviceRequestID, stepID)
	stepsStatuses, err := repository.GetDB().FindAllStepStatusByServiceRequestIDAndStepID(serviceRequestID, stepID)
	if err != nil {
		log.Errorf("No record found with given service request id %s", serviceRequestID)
		return []*models.StepsStatus{}, err
	}
	return stepsStatuses, err
}

func PrepareStepStatusResponse(srvReqID uuid.UUID, workflow *models.Workflow, stepsStatusArr []*models.StepsStatus) *models.ServiceRequestStatusResponse {
	var srvReqStatusRes models.ServiceRequestStatusResponse
	srvReqStatusRes.ServiceRequestID = srvReqID
	srvReqStatusRes.WorkflowName = workflow.Name
	stepsStatusRes := make([]models.StepStatusResponse, len(stepsStatusArr))
	stepsCount := len(workflow.Steps)
	if len(stepsStatusArr) > 0 {
		var startedStepsCount, completedStepsCount, failedStepsCount, pausedStepsCount, skippedStepsCount int
		for i, stepsStatus := range stepsStatusArr {
			stepsStatusRes[i] = models.StepStatusResponse{
				ID:        stepsStatus.StepID,
				Name:      stepsStatus.StepName,
				Status:    stepsStatus.Status,
				TimeTaken: stepsStatus.TotalTimeInMs,
				Payload:   stepsStatus.Payload,
			}
			switch stepsStatus.Status {
			case models.StatusStarted:
				startedStepsCount++
			case models.StatusCompleted:
				completedStepsCount++
			case models.StatusFailed:
				failedStepsCount++
			case models.StatusPaused:
				pausedStepsCount++
			case models.StatusSkipped:
				skippedStepsCount++
			}
		}
		// TODO Need to  change this  logic
		if startedStepsCount >= stepsCount && skippedStepsCount+completedStepsCount >= stepsCount {
			srvReqStatusRes.Status = models.StatusCompleted
		} else if failedStepsCount > 0 {
			srvReqStatusRes.Status = models.StatusFailed
		} else if pausedStepsCount > 0 {
			srvReqStatusRes.Status = models.StatusPaused
		} else if startedStepsCount != stepsCount || completedStepsCount != stepsCount {
			srvReqStatusRes.Status = models.StatusInprogress
		}
		timeTaken := calculateTimeTaken(stepsStatusArr[0].CreatedAt, stepsStatusArr[len(stepsStatusArr)-1].CreatedAt)
		srvReqStatusRes.TotalTimeInMs = timeTaken.Nanoseconds() / utils.MilliSecondsDivisor
		srvReqStatusRes.Steps = stepsStatusRes
	}
	return &srvReqStatusRes
}

func calculateTimeTaken(startTime time.Time, endTime time.Time) time.Duration {
	//log.Debug("Time Difference is == ", endTime.Sub(startTime))
	return endTime.Sub(startTime)
}
