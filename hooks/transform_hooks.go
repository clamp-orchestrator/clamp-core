package hooks

import (
	"encoding/json"

	"github.com/qntfy/kazaam"
	log "github.com/sirupsen/logrus"
)

// TransformHook : Request Transform hook
type TransformHook struct {
}

// ShouldStepExecute : Implementation method but its not required for Request Transformation
func (e *TransformHook) ShouldStepExecute(string, map[string]interface{}, string) (bool, error) {
	panic("implement me")
}

// TransformRequest : Implementation method where request will be transformed based on transformed structure
func (e *TransformHook) TransformRequest(
	stepRequestBody map[string]interface{}, transformedStructure map[string]interface{}) (map[string]interface{}, error) {
	var transformedRequestBody map[string]interface{}

	marshal, err := json.Marshal(stepRequestBody)

	specString := prepareSpecStringAsPerKazaamContract(transformedStructure, err)

	// Main kazaam transformation object
	transform, kazaamErr := kazaam.NewKazaam(string(specString))
	if kazaamErr != nil {
		// TODO If transformation fails what to do, Need to handle that scenario
		log.Errorf("Kazaam creation failed: %s", kazaamErr)
		return nil, kazaamErr
	}
	// Actual transformation happens here
	bytes, err := transform.Transform(marshal)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &transformedRequestBody)
	if err != nil {
		return nil, err
	}

	log.Debug("Evaluted value ", transformedRequestBody)
	return transformedRequestBody, nil
}

func prepareSpecStringAsPerKazaamContract(transformedStructure map[string]interface{}, _ error) []byte {
	// Operation is set to Shift mode of Kazaam
	spec := make([]map[string]interface{}, 1)
	specInterface := map[string]interface{}{
		"operation": "shift",
		"spec":      transformedStructure,
	}
	spec[0] = specInterface
	specString, err := json.Marshal(spec)
	if err != nil {
		log.Errorf("error while marshaling spec: %s", err)
	}

	return specString
}

func GetTransformHook() Hook {
	return &TransformHook{}
}
