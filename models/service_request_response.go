package models

import "github.com/google/uuid"

// ServiceRequestResponse is a structure to store the service request details
type ServiceRequestResponse struct {
	URL    string    `json:"pollUrl"`
	Status Status    `json:"status"`
	ID     uuid.UUID `json:"id"`
}
