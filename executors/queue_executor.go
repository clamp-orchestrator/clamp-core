package executors

import "errors"

type QueueVal struct {
	URL    string `json:"url" binding:"required,url"`
	Topics string `json:"headers" binding:"required"`
}

func (queueVal QueueVal) DoExecute() (interface{}, error) {
	return errors.New("queue"), errors.New("queue")
}
