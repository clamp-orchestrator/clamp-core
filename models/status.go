package models

// A Status represents the status of the service request
type Status string

const (
	// STATUS_NEW represents the step is just created
	STATUS_NEW Status = "NEW"

	// STATUS_STARTED represents the step is started execution
	STATUS_STARTED Status = "STARTED"

	// STATUS_RESUMED represents the step execution is resumed
	STATUS_RESUMED Status = "RESUMED"

	// STATUS_PAUSED represents the step execution is paused
	STATUS_PAUSED Status = "PAUSED"

	// STATUS_COMPLETED represents the step execution is completed
	STATUS_COMPLETED Status = "COMPLETED"

	// STATUS_FAILED represents the step execution is failed
	STATUS_FAILED Status = "FAILED"

	// STATUS_INPROGRESS represents the step execution is in-progress
	STATUS_INPROGRESS Status = "IN_PROGRESS"

	// STATUS_SKIPPED represents the step execution is skipped
	STATUS_SKIPPED Status = "SKIPPED"
)
