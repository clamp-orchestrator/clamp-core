package services

import (
	"clamp-core/models"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const ServiceRequestChannelSize = 1000
const ServiceRequestWorkersSize = 100

var (
	serviceRequestChannel chan models.ServiceRequest
	singletonOnce         sync.Once
)

func createServiceRequestChannel() chan models.ServiceRequest {
	singletonOnce.Do(func() {
		serviceRequestChannel = make(chan models.ServiceRequest, ServiceRequestChannelSize)
	})
	return serviceRequestChannel
}

func init() {
	createServiceRequestChannel()
	createServiceRequestWorkers()
}

func createServiceRequestWorkers() {
	for i := 0; i < ServiceRequestWorkersSize; i++ {
		go worker(i, serviceRequestChannel)
	}
}

func worker(workerId int, serviceReqChan <-chan models.ServiceRequest) {
	prefix := fmt.Sprintf("[WORKER_%d] : ", workerId)
	prefix = fmt.Sprintf("%15s", prefix)
	fmt.Printf("%s Started listening to service request channel\n", prefix)
	var stepStatus models.StepsStatus
	for serviceReq := range serviceReqChan {
		stepStatus.ServiceRequestId = serviceReq.ID
		stepStatus.WorkflowName = serviceReq.WorkflowName

		start := time.Now()
		fmt.Printf("%s Started processing service request id %s\n", prefix, serviceReq.ID)
		workflow, err := FindWorkflowByName(serviceReq.WorkflowName)
		if err == nil {
			for _, step := range workflow.Steps {
				stepStartTime := time.Now()
				stepStatus.Status =  models.STATUS_STARTED
				stepStatus.StepName = step.Name
				stepStatus.TotalTimeInMs = time.Since(stepStartTime).Milliseconds()
				SaveStepStatus(stepStatus)
				stepStartTime = time.Now()
				log.Printf("%s Started executing step id %s\n", prefix, step.Id)
				var httpClient = &http.Client{
					Timeout: time.Second * 10,
				}
				request, err := http.NewRequest(step.Mode, step.URL, nil)
				if err != nil {
					stepStatus.Status =  models.STATUS_FAILED
					stepStatus.Reason = err.Error()
					stepStatus.TotalTimeInMs = time.Since(stepStartTime).Milliseconds()
					SaveStepStatus(stepStatus)
					panic(err)
				}
				resp, err := httpClient.Do(request)
				if err != nil {
					stepStatus.Status =  models.STATUS_FAILED
					stepStatus.Reason = err.Error()
					stepStatus.TotalTimeInMs = time.Since(stepStartTime).Milliseconds()
					SaveStepStatus(stepStatus)
					panic(err)
				}
				if resp != nil {
					data, _ := ioutil.ReadAll(resp.Body)
					log.Printf("%s resp %s", prefix, string(data))
					log.Printf("%s resp %s\n", prefix, resp.Status)
					log.Printf("%s resp %d\n", prefix, resp.StatusCode)
					log.Printf("%s err %s\n", prefix, err)
					stepStatus.Status =  models.STATUS_COMPLETED
					stepStatus.TotalTimeInMs = time.Since(stepStartTime).Milliseconds()
					SaveStepStatus(stepStatus)
				}
				//stepElapsedTime := time.Since(start)
			}
		}
		elapsed := time.Since(start)
		fmt.Printf("%s Completed processing service request id %s in %s\n", prefix, serviceReq.ID, elapsed)
	}
}

func getServiceRequestChannel() chan models.ServiceRequest {
	if serviceRequestChannel == nil {
		panic(errors.New("service request channel not initialized"))
	}
	return serviceRequestChannel
}

func AddServiceRequestToChannel(serviceReq models.ServiceRequest) {
	channel := getServiceRequestChannel()
	channel <- serviceReq
}
