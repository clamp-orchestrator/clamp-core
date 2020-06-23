package services

import (
	"clamp-core/models"
	"clamp-core/repository"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log"
	"strconv"
	"time"
)

var (
	stepRequestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "create_step_request_handler_counter",
		Help: "The total number of step requests created",
	}, []string{"step_name", "time_in_ms","service_request_id"})
	stepRequestHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "create_step_request_handler_histogram",
		Help: "The total number of step requests created",
	})
)
func SaveStepStatus(stepStatusReq models.StepsStatus) (models.StepsStatus, error) {
	log.Printf("Saving step status : %v", stepStatusReq)
	stepStatusReq, err := repository.GetDB().SaveStepStatus(stepStatusReq)
	if err != nil {
		log.Printf("Failed saving step status : %v, %s", stepStatusReq, err.Error())
	}
	if stepStatusReq.Status == models.STATUS_COMPLETED{
		stepRequestCounter.WithLabelValues(stepStatusReq.StepName,strconv.FormatInt(stepStatusReq.TotalTimeInMs, 10),stepStatusReq.ServiceRequestId.String()).Inc()
	}
	return stepStatusReq, err
}

func FindStepStatusByServiceRequestId(serviceRequestId uuid.UUID) ([]models.StepsStatus, error) {
	log.Printf("Find step statues by request id : %s ", serviceRequestId)
	stepsStatuses, err := repository.GetDB().FindStepStatusByServiceRequestId(serviceRequestId)
	if err != nil {
		log.Printf("No record found with given service request id %s", serviceRequestId)
		return []models.StepsStatus{}, err
	}
	return stepsStatuses, err
}

func FindStepStatusByServiceRequestIdAndStatus(serviceRequestId uuid.UUID, status models.Status) ([]models.StepsStatus, error) {
	log.Printf("Find step statues by request id : %s ", serviceRequestId)
	stepsStatuses, err := repository.GetDB().FindStepStatusByServiceRequestIdAndStatus(serviceRequestId, status)
	if err != nil {
		log.Printf("No record found with given service request id %s", serviceRequestId)
		return []models.StepsStatus{}, err
	}
	return stepsStatuses, err
}

func FindAllStepStatusByServiceRequestIdAndStepId(serviceRequestId uuid.UUID, stepId int) ([]models.StepsStatus, error) {
	log.Printf("Find all step statues by request id : %s and step id : %d", serviceRequestId, stepId)
	stepsStatuses, err := repository.GetDB().FindAllStepStatusByServiceRequestIdAndStepId(serviceRequestId, stepId)
	if err != nil {
		log.Printf("No record found with given service request id %s", serviceRequestId)
		return []models.StepsStatus{}, err
	}
	return stepsStatuses, err
}

func PrepareStepStatusResponse(srvReqId uuid.UUID, workflow models.Workflow, stepsStatusArr []models.StepsStatus) models.ServiceRequestStatusResponse {
	var srvReqStatusRes models.ServiceRequestStatusResponse
	srvReqStatusRes.ServiceRequestId = srvReqId
	srvReqStatusRes.WorkflowName = workflow.Name
	stepsStatusRes := make([]models.StepStatusResponse, len(stepsStatusArr))
	stepsCount := len(workflow.Steps)
	if len(stepsStatusArr) > 0 {
		var startedStepsCount, completedStepsCount, failedStepsCount, pausedStepsCount, skippedStepsCount int
		for i, stepsStatus := range stepsStatusArr {
			stepsStatusRes[i] = models.StepStatusResponse{
				Id:        stepsStatus.StepId,
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
		srvReqStatusRes.TotalTimeInMs = timeTaken.Nanoseconds() / models.MilliSecondsDivisor
		srvReqStatusRes.Steps = stepsStatusRes
	}
	return srvReqStatusRes
}

func calculateTimeTaken(startTime time.Time, endTime time.Time) time.Duration {
	//log.Println("Time Difference is == ", endTime.Sub(startTime))
	return endTime.Sub(startTime)
}
