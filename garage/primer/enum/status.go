package enum

// Statuses

type Status string

const (

	// Open
	Open Status = "OPEN"

	// Closed
	Closed Status = "CLOSED"

	// Active
	Active Status = "ACTIVE"

	// Inactive
	Inactive Status = "INACTIVE"

	// Pending
	Pending Status = "PENDING"

	// Deleted
	Deleted Status = "DELETED"

	// Completed
	Completed Status = "COMPLETED"
)

const (
	// Success
	Success Status = "SUCCESS"

	// Failure
	Failure Status = "FAILURE"
)
