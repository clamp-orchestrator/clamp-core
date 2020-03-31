package services

import (
	"clamp-core/models"
	"sync"
)

const SERVICE_REQUEST_CHANNEL_SIZE = 100

var (
	serviceRequestChannel chan models.ServiceRequest
	singletonOnce         sync.Once
)

func createServiceRequestChannel() chan models.ServiceRequest {
	singletonOnce.Do(func() {
		serviceRequestChannel = make(chan models.ServiceRequest, SERVICE_REQUEST_CHANNEL_SIZE)
	})
	return serviceRequestChannel
}
