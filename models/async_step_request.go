package models

type AsyncStepRequest struct {
	StepStatus StepsStatus            `json:"step_status"`
	Step       Step                   `json:"step"`
	Payload    map[string]interface{} `json:"payload"`
	Prefix     string                 `json:"prefix"`
}
