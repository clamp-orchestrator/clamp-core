package models

import "encoding/xml"

type XMLStepRequest struct {
	XMLName xml.Name `xml:"request"`
	Text    string   `xml:",chardata"`
	Body    struct {
		Text             string `xml:",chardata"`
		ServiceRequestID string `xml:"serviceRequestId"`
		StepID           string `xml:"stepId"`
		Payload          string `xml:"payload"`
	} `xml:"body"`
}
