package hooks

type Hook interface {
	ShouldStepExecute(string, map[string]interface{}, string) (bool, error)
	TransformRequest(map[string]interface{}, map[string]interface{}) (map[string]interface{}, error)
}

type defaultHook struct {
}

func (d defaultHook) TransformRequest(m map[string]interface{}, s map[string]interface{}) (map[string]interface{}, error) {
	return m, nil
}

func (d defaultHook) ShouldStepExecute(s string, m map[string]interface{}, s2 string) (bool, error) {
	return true, nil
}

func GetDefaultHook() Hook {
	return defaultHook{}
}
