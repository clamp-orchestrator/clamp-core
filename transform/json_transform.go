package transform

import (
	"clamp-core/hooks"

	log "github.com/sirupsen/logrus"
)

type JSONTransform struct {
	Spec map[string]interface{} `json:"spec"`
}

func (jsonTransform JSONTransform) DoTransform(requestBody map[string]interface{}, prefix string) (map[string]interface{}, error) {
	log.Debugf("%s Json Transformation : Transform keys %v and request body:%v", prefix, jsonTransform.Spec, requestBody)
	transformedRequestBody, err := hooks.GetTransformHook().TransformRequest(requestBody, jsonTransform.Spec)
	if err != nil {
		log.Debugf("Transformation failed: %s", err)
		return nil, err
	}
	log.Debug("Transformed Request Body is ", transformedRequestBody)
	return transformedRequestBody, nil
}
