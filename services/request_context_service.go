package services

import "clamp-core/models"

func CreateRequestContext(workflow models.Workflow, request models.ServiceRequest) (context models.RequestContext) {
	context = models.RequestContext{
		ServiceRequestId: request.ID,
		WorkflowName:     workflow.Name,
	}
	context.StepsContext = make(map[string]*models.StepContext)
	for _, step := range workflow.Steps {
		context.StepsContext[step.Name] = &models.StepContext{
			Request:  nil,
			Response: nil,
		}
	}
	return
}

func EnhanceRequestContextWithExecutedSteps(context *models.RequestContext) {
	stepsStatuses, err := FindStepStatusByServiceRequestIdAndStatusOrderByCreatedAtDesc(context.ServiceRequestId, models.STATUS_COMPLETED)
	if err == nil {
		for _, stepsStatus := range stepsStatuses {
			context.StepsContext[stepsStatus.StepName] = &models.StepContext{
				Request:  stepsStatus.Payload.Request,
				Response: stepsStatus.Payload.Response,
			}
		}
	}
}
