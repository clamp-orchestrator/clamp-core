package domain

//Workflow is a structure to store the service request details
type Request struct {
	ServiceFlow ServiceFlow
}
type ServiceFlow struct {
	Description string `json:"description"`
	FlowMode    string `json:"flowMode"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Steps       Steps  `json:"steps"`
}

type Steps struct {
	Step []Step `json:"step"`
}

type Step struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

//Create a new work flow for a given service flow and return service flow details
func CreateWorkflow(serviceFlow ServiceFlow) ServiceFlow {
	return newServiceFlow(serviceFlow)
}

func newServiceFlow(serviceFlow ServiceFlow) ServiceFlow {
	return ServiceFlow{Description: serviceFlow.Description, FlowMode: serviceFlow.FlowMode, Id: serviceFlow.Id, Name: serviceFlow.Name, Enabled: serviceFlow.Enabled, Steps: serviceFlow.Steps}
}
