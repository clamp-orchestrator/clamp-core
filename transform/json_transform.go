package transform

import (
	"clamp-core/hooks"
	"log"
)

type JSONTransform struct {
	Spec map[string]interface{} `json:"spec"`
}

func (jsonTransform JSONTransform) DoTransform(requestBody map[string]interface{}, prefix string) (map[string]interface{}, error) {
	log.Printf("%s Json Transformation : Transform keys %v and request body:%v", prefix, jsonTransform.Spec, requestBody)
	transformedRequestBody, err := hooks.GetTransformHook().TransformRequest(requestBody, jsonTransform.Spec)
	if err != nil {
		log.Println("Transformation failed")
		return nil, err
	}
	log.Println("Transformed Request Body is ", transformedRequestBody)
	return transformedRequestBody, nil
}
