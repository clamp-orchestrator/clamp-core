package models

import "github.com/google/uuid"

//ServiceRequest is a structure to store the service request details
type ServiceRequestResponse struct {
	URL    string    `json:"pollUrl"`
	Status Status    `json:"status"`
	ID     uuid.UUID `json:"serviceRequestId"`
}
