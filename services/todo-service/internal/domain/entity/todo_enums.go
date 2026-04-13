package entity

type TodoStatus string

const (
	TodoStatusPending    TodoStatus = "PENDING"
	TodoStatusInProgress TodoStatus = "IN_PROGRESS"
	TodoStatusDone       TodoStatus = "DONE"
)

type Priority string

const (
	PriorityLow    Priority = "LOW"
	PriorityMedium Priority = "MEDIUM"
	PriorityHigh   Priority = "HIGH"
	PriorityUrgent Priority = "URGENT"
)
