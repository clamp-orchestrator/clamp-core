package executors

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// HTTPVal : Http configuration details
type HTTPVal struct {
	Method  string `json:"method" binding:"required,oneof=GET POST PUT PATCH DELETE"`
	URL     string `json:"url" binding:"required,url"`
	Headers string `json:"headers"`
}

// DoExecute : Preparing to make a http call with request body
func (httpVal *HTTPVal) DoExecute(requestBody interface{}, prefix string) (interface{}, error) {
	log.Debugf("%s HTTP Executor: Calling http %s:%s body:%v", prefix, httpVal.Method, httpVal.URL, requestBody)
	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
	requestJSONBytes, _ := json.Marshal(requestBody)
	request, err := http.NewRequest(httpVal.Method, httpVal.URL, bytes.NewBuffer(requestJSONBytes))
	fetchAndLoadRequestWithHeadersIfDefined(httpVal, request)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		data, _ := ioutil.ReadAll(resp.Body)
		err = errors.New(string(data))
		log.Error("Unable to execute \t", httpVal.URL, " with error message", err)
		return nil, err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("%sHTTP Executor: Successfully called http %s:%s", prefix, httpVal.Method, httpVal.URL)
	return string(data), err
}

func fetchAndLoadRequestWithHeadersIfDefined(httpVal *HTTPVal, request *http.Request) {
	if httpVal.Headers != "" {
		httpHeaders := strings.Split(httpVal.Headers, ";")
		for _, header := range httpHeaders[:len(httpHeaders)-1] {
			httpHeader := strings.Split(header, ":")
			if len(httpHeader) > 1 {
				request.Header.Add(httpHeader[0], httpHeader[1])
			}
		}
	}
}
