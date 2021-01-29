package services

import "clamp-core/models"

func CreateRequestContext(workflow models.Workflow, request models.ServiceRequest) (context models.RequestContext) {
	context = models.RequestContext{
		ServiceRequestID: request.ID,
		WorkflowName:     workflow.Name,
	}
	context.StepsContext = make(map[string]*models.StepContext)
	for _, step := range workflow.Steps {
		context.StepsContext[step.Name] = &models.StepContext{
			Request:         nil,
			Response:        nil,
			RequestHeaders:  request.RequestHeaders,
			ResponseHeaders: "",
		}
	}
	return
}

func EnhanceRequestContextWithExecutedSteps(context *models.RequestContext) {
	stepsStatuses, err := FindStepStatusByServiceRequestIDAndStatus(context.ServiceRequestID, models.STATUS_COMPLETED)
	if err == nil {
		for _, stepsStatus := range stepsStatuses {
			context.StepsContext[stepsStatus.StepName] = &models.StepContext{
				Request:  stepsStatus.Payload.Request,
				Response: stepsStatus.Payload.Response,
			}
		}
	}
}

func ComputeRequestToCurrentStepInContext(workflow models.Workflow, currentStepExecuting models.Step, requestContext *models.RequestContext, stepIndex int, stepRequestPayload map[string]interface{}) {
	if requestContext.GetStepRequestFromContext(currentStepExecuting.Name) == nil {
		if stepIndex == 0 {
			//for first step in execution
			requestContext.SetStepRequestToContext(currentStepExecuting.Name, stepRequestPayload)
		}
		if stepIndex > 0 {
			prevStepExecuted := workflow.Steps[stepIndex-1]
			prevStepResponse := requestContext.GetStepResponseFromContext(prevStepExecuted.Name)
			if prevStepResponse != nil {
				requestContext.SetStepRequestToContext(currentStepExecuting.Name, prevStepResponse)
			} else {
				//for skipped step there will be no response
				prevStepRequest := requestContext.GetStepRequestFromContext(prevStepExecuted.Name)
				requestContext.SetStepRequestToContext(currentStepExecuting.Name, prevStepRequest)
			}
		}
	}
}
