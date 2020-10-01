package hooks

import (
	"encoding/json"
	"log"

	"github.com/qntfy/kazaam"
)

// TransformHook : Request Transform hook
type TransformHook struct {
}

// ShouldStepExecute : Implementation method but its not required for Request Transformation
func (e *TransformHook) ShouldStepExecute(string, map[string]interface{}, string) (bool, error) {
	panic("implement me")
}

// TransformRequest : Implementation method where request will be transformed based on transformed structure
func (e *TransformHook) TransformRequest(stepRequestBody map[string]interface{}, transformedStructure map[string]interface{}) (map[string]interface{}, error) {
	var transformedRequestBody map[string]interface{}

	marshal, err := json.Marshal(stepRequestBody)

	specString := prepareSpecStringAsPerKazaamContract(transformedStructure, err)

	//Main kazaam transformation object
	transform, kazaamErr := kazaam.NewKazaam(string(specString))
	if kazaamErr != nil {
		//TODO If transformation fails what to do, Need to handle that scenario
		log.Println("Something went wrong")
		return nil, kazaamErr
	}
	//Actual transformation happens here
	bytes, err := transform.Transform(marshal)

	_ = json.Unmarshal(bytes, &transformedRequestBody)
	log.Println("Evaluted value ", transformedRequestBody)
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
