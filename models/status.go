package models

type Status string

const (
	STATUS_NEW        Status = "NEW"
	STATUS_STARTED    Status = "STARTED"
	STATUS_RESUMED    Status = "RESUMED"
	STATUS_PAUSED     Status = "PAUSED"
	STATUS_COMPLETED  Status = "COMPLETED"
	STATUS_FAILED     Status = "FAILED"
	STATUS_INPROGRESS Status = "IN_PROGRESS"
	STATUS_SKIPPED    Status = "STATUS_SKIPPED"
)
