package hooks

type hookInterface interface {
	preStepExecution(string, map[string]interface{}, string) (bool, error)
}

var hook hookInterface

//TODO: need to refactor to support multiple hooks
func PreStepHookExecutor(whenCondition string, stepRequest map[string]interface{}, prefix string) (bool, error) {
	if whenCondition != "" {
		hook = &ExprHook{}
		return hook.preStepExecution(whenCondition, stepRequest, prefix)
	}
	panic("no implementation")
}
