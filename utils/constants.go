package utils

const (
	StepTypeAsync = "ASYNC"
	StepTypeSync  = "SYNC"
)

const ServiceRequestChannelSize = 1000
const ServiceRequestWorkersSize = 100

const ResumeStepResponseChannelSize = 1000
const ResumeStepResponseWorkersSize = 100

var MilliSecondsDivisor int64 = 1000000

const (
	StepModeHTTP  = "HTTP"
	StepModeKafka = "KAFKA"
	StepModeAMQP  = "AMQP"
)

const (
	TransformFormatXML  = "XML"
	TransformFormatJSON = "JSON"
)
