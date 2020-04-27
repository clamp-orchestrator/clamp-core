package models

import (
	"fmt"
	"testing"
)


func Test(t *testing.T) {
	var employee = make(map[string]RequestResponse)
	request := RequestResponse{
		Request:  prepareRequestPayload(),
		Response: prepareRequestPayload(),
	}
	employee["Mark"] = request
	employee["Sandy"] = request
	fmt.Println(employee)
}