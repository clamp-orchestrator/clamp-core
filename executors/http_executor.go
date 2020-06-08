package executors

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type HttpVal struct {
	Method  string `json:"method" binding:"required,oneof=GET POST PUT PATCH DELETE"`
	Url     string `json:"url" binding:"required,url"`
	Headers string `json:"headers"`
}

func (httpVal HttpVal) DoExecute(requestBody interface{}, prefix string) (interface{}, error) {
	log.Printf("%s HTTP Executor: Calling http %s:%s body:%v", prefix, httpVal.Method, httpVal.Url, requestBody)
	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
	requestJsonBytes, _ := json.Marshal(requestBody)
	request, err := http.NewRequest(httpVal.Method, httpVal.Url, bytes.NewBuffer(requestJsonBytes))
	fetchAndLoadRequestWithHeadersIfDefined(httpVal, request)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		data, _ := ioutil.ReadAll(resp.Body)
		err := errors.New(string(data))
		log.Println("Unable to execute \t", httpVal.Url, " with error message", err)
		return nil, err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	log.Printf("%sHTTP Executor: Successfully called http %s:%s", prefix, httpVal.Method, httpVal.Url)
	return string(data), err
}

func fetchAndLoadRequestWithHeadersIfDefined(httpVal HttpVal, request *http.Request) {
	if httpVal.Headers != "" {
		httpHeaders := strings.Split(httpVal.Headers, ";")
		for _, header := range httpHeaders[:len(httpHeaders)-1] {
			httpHeader := strings.Split(header, ":")
			request.Header.Add(httpHeader[0], httpHeader[1])
		}
	}
}
