package event

import "time"

type TodoListCreated struct {
	TodoListID int64     `json:"todo_list_id"`
	ActorID    int64     `json:"actor_id"`
	Name       string    `json:"name"`
	OccurredOn time.Time `json:"occurred_at"`
}

func (e TodoListCreated) EventName() string {
	return "todo_list.created"
}

func (e TodoListCreated) OccurredAt() time.Time {
	return e.OccurredOn
}

type TodoListUpdated struct {
	TodoListID int64     `json:"todo_list_id"`
	ActorID    int64     `json:"actor_id"`
	OccurredOn time.Time `json:"occurred_at"`
}

func (e TodoListUpdated) EventName() string {
	return "todo_list.updated"
}

func (e TodoListUpdated) OccurredAt() time.Time {
	return e.OccurredOn
}

type TodoListDeleted struct {
	TodoListID int64     `json:"todo_list_id"`
	ActorID    int64     `json:"actor_id"`
	OccurredOn time.Time `json:"occurred_at"`
}

func (e TodoListDeleted) EventName() string {
	return "todo_list.deleted"
}

func (e TodoListDeleted) OccurredAt() time.Time {
	return e.OccurredOn
}
