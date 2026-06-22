package event

import "time"

// TodoAssigned is emitted when a todo is assigned to a new user.
// Only emitted when assignee_id actually changes.
type TodoAssigned struct {
	TodoID     int64     `json:"todo_id"`
	TodoListID int64     `json:"todo_list_id"`
	ActorID    int64     `json:"actor_id"`    // who performed the assignment
	AssigneeID int64     `json:"assignee_id"` // who was assigned
	Title      string    `json:"title"`
	OccurredOn time.Time `json:"occurred_at"`
}

func (e TodoAssigned) EventName() string     { return "todo.assigned" }
func (e TodoAssigned) OccurredAt() time.Time { return e.OccurredOn }
