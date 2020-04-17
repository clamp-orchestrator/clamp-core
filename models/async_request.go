package models

type AsyncStepExecutionRequest struct {
	StepStatus StepsStatus            `json:"stepStatus"`
	Step       Step                   `json:"step"`
	Payload    map[string]interface{} `json:"payload"`
	Prefix     string                 `json:"prefix"`
}
