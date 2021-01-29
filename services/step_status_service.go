package services

import (
	"clamp-core/models"
	"clamp-core/repository"
	"clamp-core/utils"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

func SaveStepStatus(stepStatusReq models.StepsStatus) (models.StepsStatus, error) {
	log.Printf("Saving step status : %v", stepStatusReq)
	stepStatusReq, err := repository.GetDB().SaveStepStatus(stepStatusReq)
	if err != nil {
		log.Printf("Failed saving step status : %v, %s", stepStatusReq, err.Error())
	}
	serviceRequestStepNameTimeExecutorCounter.WithLabelValues(stepStatusReq.ServiceRequestID.String(), stepStatusReq.StepName, string(stepStatusReq.Status)).Add(float64(stepStatusReq.TotalTimeInMs))
	serviceRequestStepNameCounter.WithLabelValues(stepStatusReq.ServiceRequestID.String(), stepStatusReq.StepName, string(stepStatusReq.Status)).Inc()
	serviceRequestStepNameTimeExecutorHistogram.Observe(float64(stepStatusReq.TotalTimeInMs))
	return stepStatusReq, err
}

func FindStepStatusByServiceRequestID(serviceRequestID uuid.UUID) ([]models.StepsStatus, error) {
	log.Printf("Find step statues by request id : %s ", serviceRequestID)
	stepsStatuses, err := repository.GetDB().FindStepStatusByServiceRequestID(serviceRequestID)
	if err != nil {
		log.Printf("No record found with given service request id %s", serviceRequestID)
		return []models.StepsStatus{}, err
	}
	return stepsStatuses, err
}

func FindStepStatusByServiceRequestIDAndStatus(serviceRequestID uuid.UUID, status models.Status) ([]models.StepsStatus, error) {
	log.Printf("Find step statues by request id : %s ", serviceRequestID)
	stepsStatuses, err := repository.GetDB().FindStepStatusByServiceRequestIDAndStatus(serviceRequestID, status)
	if err != nil {
		log.Printf("No record found with given service request id %s", serviceRequestID)
		return []models.StepsStatus{}, err
	}
	return stepsStatuses, err
}

func FindAllStepStatusByServiceRequestIDAndStepID(serviceRequestID uuid.UUID, stepID int) ([]models.StepsStatus, error) {
	log.Printf("Find all step statues by request id : %s and step id : %d", serviceRequestID, stepID)
	stepsStatuses, err := repository.GetDB().FindAllStepStatusByServiceRequestIDAndStepID(serviceRequestID, stepID)
	if err != nil {
		log.Printf("No record found with given service request id %s", serviceRequestID)
		return []models.StepsStatus{}, err
	}
	return stepsStatuses, err
}

func PrepareStepStatusResponse(srvReqID uuid.UUID, workflow models.Workflow, stepsStatusArr []models.StepsStatus) models.ServiceRequestStatusResponse {
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
			case models.STATUS_STARTED:
				startedStepsCount++
			case models.STATUS_COMPLETED:
				completedStepsCount++
			case models.STATUS_FAILED:
				failedStepsCount++
			case models.STATUS_PAUSED:
				pausedStepsCount++
			case models.STATUS_SKIPPED:
				skippedStepsCount++
			}
		}
		//TODO Need to  change this  logic
		if startedStepsCount >= stepsCount && skippedStepsCount+completedStepsCount >= stepsCount {
			srvReqStatusRes.Status = models.STATUS_COMPLETED
		} else if failedStepsCount > 0 {
			srvReqStatusRes.Status = models.STATUS_FAILED
		} else if pausedStepsCount > 0 {
			srvReqStatusRes.Status = models.STATUS_PAUSED
		} else if startedStepsCount != stepsCount || completedStepsCount != stepsCount {
			srvReqStatusRes.Status = models.STATUS_INPROGRESS
		}
		timeTaken := calculateTimeTaken(stepsStatusArr[0].CreatedAt, stepsStatusArr[len(stepsStatusArr)-1].CreatedAt)
		srvReqStatusRes.TotalTimeInMs = timeTaken.Nanoseconds() / utils.MilliSecondsDivisor
		srvReqStatusRes.Steps = stepsStatusRes
	}
	return srvReqStatusRes
}

func calculateTimeTaken(startTime time.Time, endTime time.Time) time.Duration {
	//log.Println("Time Difference is == ", endTime.Sub(startTime))
	return endTime.Sub(startTime)
}
