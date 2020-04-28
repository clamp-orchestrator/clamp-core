package transform

type HtmlTransform struct {
	//TODO Capture additional details which is required for transformation
	Keys	map[string]interface{} `json:"keys"`
}

func (t HtmlTransform) DoTransform(requestBody map[string]interface{}, prefix string) (map[string]interface{}, error){
	return requestBody, nil
}