package transform

import (
	"clamp-core/hooks"
	"fmt"
	"log"
)

type JsonTransform struct {
	Keys	map[string]interface{} `json:"keys"`
}

func (jsonTransform JsonTransform) DoTransform(requestBody map[string]interface{}, prefix string) (map[string]interface{}, error){
	log.Printf("%s Json Transformation : Transform keys %v and request body:%v", prefix, jsonTransform.Keys, requestBody)
	transformedRequestBody := make(map[string]interface{})
	for key, requestBodyKey := range jsonTransform.Keys {
		jsonPathKey := fmt.Sprintf("%v", requestBodyKey)
		jsonPathValue, err := hooks.GetTransformHook().TransformRequest(requestBody, jsonPathKey)
		if err != nil  {
			log.Println("Transformation failed")
			return nil, err
		}
		transformedRequestBody[key] = jsonPathValue[jsonPathKey]
	}
	log.Println("Transformed Request Body is ", transformedRequestBody)
	return transformedRequestBody, nil
}