package models

import (
	"github.com/google/uuid"
)

type AsyncStepResponse struct {
	ServiceRequestId   uuid.UUID              `json:"serviceRequestId"`
	StepId             int                    `json:"stepId"`
	Response           map[string]interface{} `json:"response"`
	Error              ClampErrorResponse     `json:"error"`
	stepStatusRecorded bool
}

func (res *AsyncStepResponse) SetStepStatusRecorded(stepStatusRecorded bool) {
	res.stepStatusRecorded = stepStatusRecorded
}

func (res *AsyncStepResponse) IsStepStatusRecorded() bool {
	return res.stepStatusRecorded
}

//StepProcessed to be true if step response was received internally, false if response received from external
