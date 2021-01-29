package models

import (
	"github.com/google/uuid"
)

type AsyncStepResponse struct {
	ServiceRequestID   uuid.UUID              `json:"serviceRequestId"`
	StepID             int                    `json:"stepId"`
	Response           map[string]interface{} `json:"response"`
	Error              ClampErrorResponse     `json:"error"`
	stepStatusRecorded bool
	RequestHeaders     string
}

func (res *AsyncStepResponse) SetStepStatusRecorded(stepStatusRecorded bool) {
	res.stepStatusRecorded = stepStatusRecorded
}

func (res *AsyncStepResponse) IsStepStatusRecorded() bool {
	return res.stepStatusRecorded
}

//StepProcessed to be true if step response was received internally, false if response received from external
