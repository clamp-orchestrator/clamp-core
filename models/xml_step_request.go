package models

import "encoding/xml"

type XmlStepRequest struct {
	XMLName xml.Name `xml:"request"`
	Text    string   `xml:",chardata"`
	Body    struct {
		Text             string `xml:",chardata"`
		ServiceRequestId string `xml:"serviceRequestId"`
		StepId           string `xml:"stepId"`
		Payload          string `xml:"payload"`
	} `xml:"body"`
}
