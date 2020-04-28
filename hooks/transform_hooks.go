package hooks

import (
	"encoding/json"
	"github.com/qntfy/kazaam"
	"log"
)

type TransformHook struct {
}

func (e *TransformHook) ShouldStepExecute(string, map[string]interface{}, string) (bool, error) {
	panic("implement me")
}

func (e *TransformHook) TransformRequest(stepRequestBody map[string]interface{}, transformedStructure map[string]interface{}) (map[string]interface{}, error) {
	var transformedRequestBody map[string]interface{}

	marshal, err := json.Marshal(stepRequestBody)

	specString := prepareSpecStringAsPerKazaamContract(transformedStructure, err)

	//Main kazaam transformation object
	transform, kazaamErr := kazaam.NewKazaam(string(specString))
	if kazaamErr != nil {
		log.Println("Something went wrong")
	}
	//Actual transformation happens here
	bytes, err := transform.Transform(marshal)

	err = json.Unmarshal(bytes, &transformedRequestBody)
	log.Println("Evaluted value ", transformedRequestBody)
	if err != nil {
		log.Println("Evaluted error ", err)
		return nil, err
	}
	return transformedRequestBody, nil
}

func prepareSpecStringAsPerKazaamContract(transformedStructure map[string]interface{}, err error) []byte {
	//Operation is set to Shift mode of Kazaam
	spec := make([]map[string]interface{}, 1)
	specInterface := map[string]interface{}{
		"operation": "shift",
		"spec":      transformedStructure,
	}
	spec[0] = specInterface
	specString, err := json.Marshal(spec)
	return specString
}

func GetTransformHook() Hook {
	return &TransformHook{}
}
