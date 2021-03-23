package executors

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	requestJSONBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error while marshaling http executor request body: %w", err)
	}

	request, err := http.NewRequest(httpVal.Method, httpVal.URL, bytes.NewBuffer(requestJSONBytes))
	if err != nil {
		return nil, fmt.Errorf("error while creating http executor request: %w", err)
	}

	populateRequestHeaders(httpVal.Headers, &request.Header)

	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error while executing http request: %w", err)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading http executor response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("HTTP %s request failed on %s: %s", httpVal.Method, httpVal.URL, data)
	}

	log.Debugf("%sHTTP Executor: Successfully called http %s:%s", prefix, httpVal.Method, httpVal.URL)
	return string(data), err
}

func populateRequestHeaders(httpValHeaders string, header *http.Header) {
	if httpValHeaders != "" {
		for _, httpValHeader := range strings.Split(httpValHeaders, ";") {
			httpValHeaderParts := strings.Split(httpValHeader, ":")
			if len(httpValHeaderParts) > 1 {
				header.Add(httpValHeaderParts[0], httpValHeaderParts[1])
			}
		}
	}
}
