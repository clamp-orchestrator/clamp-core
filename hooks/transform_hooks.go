package hooks

type TransformHook struct {
}

func (e *TransformHook) ShouldStepExecute(string, map[string]interface{}, string) (bool, error) {
	panic("implement me")
}

func (e *TransformHook) TransformRequest(m map[string]interface{}, transformFormat string) (map[string]interface{}, error) {

	return m, nil
}

func GetTransformHook() Hook {
	return &TransformHook{}
}
