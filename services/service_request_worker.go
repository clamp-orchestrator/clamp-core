package services

import (
	"clamp-core/models"
	"errors"
	"fmt"
	"io/ioutil"
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
	for serviceReq := range serviceReqChan {
		start := time.Now()
		fmt.Printf("%s Started processing service request id %s\n", prefix, serviceReq.ID)
		workflow, err := FindWorkflowByName(serviceReq.WorkflowName)
		if err == nil {
			for _, step := range workflow.Steps {
				fmt.Printf("%s Started executing step id %s\n", prefix, step.Id)
				var httpClient = &http.Client{
					Timeout: time.Second * 10,
				}
				request, err := http.NewRequest(step.Mode, step.URL, nil)
				if err != nil {
					panic(err)
				}
				resp, err := httpClient.Do(request)
				if err != nil {
					panic(err)
				}
				if resp != nil {
					data, _ := ioutil.ReadAll(resp.Body)
					fmt.Printf("%s resp %s", prefix, string(data))
					fmt.Printf("%s resp %s\n", prefix, resp.Status)
					fmt.Printf("%s resp %d\n", prefix, resp.StatusCode)
					fmt.Printf("%s err %s\n", prefix, err)
				}
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
