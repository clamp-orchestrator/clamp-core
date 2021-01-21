package transform

type XMLTransform struct {
	// TODO Capture additional details which is required for transformation
	Keys map[string]interface{} `json:"keys"`
}

func (t XMLTransform) DoTransform(requestBody map[string]interface{}, prefix string) (map[string]interface{}, error) {
	return requestBody, nil
}
