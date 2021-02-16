package models

type Status string

const (
	StatusNew        Status = "NEW"
	StatusStarted    Status = "STARTED"
	StatusResumed    Status = "RESUMED"
	StatusPaused     Status = "PAUSED"
	StatusCompleted  Status = "COMPLETED"
	StatusFailed     Status = "FAILED"
	StatusInprogress Status = "IN_PROGRESS"
	StatusSkipped    Status = "SKIPPED"
)
