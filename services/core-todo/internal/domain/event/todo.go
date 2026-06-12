package event

import "time"

type TodoCreated struct {
	TodoID     int64     `json:"todo_id"`
	TodoListID int64     `json:"todo_list_id"`
	ActorID    int64     `json:"actor_id"`
	Title      string    `json:"title"`
	OccurredOn time.Time `json:"occurred_at"`
}

func (e TodoCreated) EventName() string {
	return "todo.created"
}

func (e TodoCreated) OccurredAt() time.Time {
	return e.OccurredOn
}

type TodoUpdated struct {
	TodoID     int64     `json:"todo_id"`
	TodoListID int64     `json:"todo_list_id"`
	ActorID    int64     `json:"actor_id"`
	OccurredOn time.Time `json:"occurred_at"`
}

func (e TodoUpdated) EventName() string {
	return "todo.updated"
}

func (e TodoUpdated) OccurredAt() time.Time {
	return e.OccurredOn
}

type TodoDeleted struct {
	TodoID     int64     `json:"todo_id"`
	TodoListID int64     `json:"todo_list_id"`
	ActorID    int64     `json:"actor_id"`
	OccurredOn time.Time `json:"occurred_at"`
}

func (e TodoDeleted) EventName() string {
	return "todo.deleted"
}

func (e TodoDeleted) OccurredAt() time.Time {
	return e.OccurredOn
}
