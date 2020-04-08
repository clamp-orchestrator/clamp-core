package executors

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type HttpVal struct {
	Method  string `json:"method" binding:"required,oneof=GET POST PUT PATCH DELETE"`
	Url     string `json:"url" binding:"required,url"`
	Headers string `json:"headers"`
}

func (httpVal HttpVal) DoExecute() (interface{}, error) {
	log.Println("Inside HTTP Do Execute function")
	prefix := log.Prefix()
	log.SetPrefix("")
	log.Printf("%s HTTP Executor: Calling http %s:%s", prefix, httpVal.Method, httpVal.Url)
	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
	request, err := http.NewRequest(httpVal.Method, httpVal.Url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		err := errors.New(string(data))
		log.Println("Unable to execute \t", httpVal.Url," with error message",err)
		return nil, err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	log.Printf("%sHTTP Executor: Successfully called http %s:%s", prefix, httpVal.Method, httpVal.Url)
	return string(data), err
}
