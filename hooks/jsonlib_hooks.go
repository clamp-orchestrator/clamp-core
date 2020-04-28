package hooks

//Example to implement new lib
type jsonLib struct {
}

func (d jsonLib) TransformRequest(m map[string]interface{}, s map[string]interface{}) (map[string]interface{}, error) {
	return m, nil
}

func (d jsonLib) ShouldStepExecute(s string, m map[string]interface{}, s2 string) (bool, error) {
	return true, nil
}

func GetJsonLib() Hook {
	return jsonLib{}
}
