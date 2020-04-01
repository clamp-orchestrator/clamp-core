package services

import (
	"clamp-core/models"
	"errors"
	"fmt"
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
	fmt.Printf("ServiceRequestWorker id %d started \n", workerId)
	for serviceReq := range serviceReqChan {
		fmt.Printf("ServiceRequestWorker id %d started processing service request id %s\n", workerId, serviceReq.ID)
		workflow, err := FindWorkflowByName(serviceReq.WorkflowName)
		if err == nil {
			for _, step := range workflow.Steps {
				fmt.Printf("Started executing step id %s", step.Id)
			}
		}
		time.Sleep(time.Second)
		fmt.Printf("ServiceRequestWorker id %d completed processing service request id %s\n", workerId, serviceReq.ID)
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
