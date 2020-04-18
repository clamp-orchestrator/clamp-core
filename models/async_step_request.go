package models

type AsyncStepRequest struct {
	StepStatus StepsStatus            `json:"stepStatus"`
	Step       Step                   `json:"step"`
	Payload    map[string]interface{} `json:"payload"`
	Prefix     string                 `json:"prefix"`
}
