package input

import "time"

type CreateNotification struct {
	ReceiverID   int64
	ResourceType string // "todo", "todo_list", ...
	ResourceID   int64
	EventName    string // "todo.assigned", ...
	OccurredAt   time.Time
	Title        string
	Content      string
}
