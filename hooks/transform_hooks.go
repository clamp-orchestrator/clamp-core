package hooks

import (
	"github.com/antonmedv/expr"
	"log"
)

type TransformHook struct {
}

func (e *TransformHook) ShouldStepExecute(string, map[string]interface{}, string) (bool, error) {
	panic("implement me")
}

func (e *TransformHook) TransformRequest(stepRequestBody map[string]interface{}, key string) (map[string]interface{}, error) {
	var transformedRequestBody map[string]interface{}
	eval, err := expr.Eval(key, stepRequestBody)
	log.Println("Evaluted value ", eval)
	log.Println("Evaluted error ", err)
	transformedRequestBody = map[string]interface{}{key:eval}
	return transformedRequestBody, nil
}

func GetTransformHook() Hook {
	return &TransformHook{}
}
