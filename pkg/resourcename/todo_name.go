package resourcename

import (
	"fmt"
	"regexp"
	"strconv"
)

// users/{user_id}/todo-lists/{list_id}
type TodoListResourceName struct {
	UserID     int64
	TodoListID int64
}

func (n TodoListResourceName) String() string {
	return fmt.Sprintf("users/%d/todo-lists/%d", n.UserID, n.TodoListID)
}

var todoListRegex = regexp.MustCompile(`^users/(\d+)/todo-lists/(\d+)$`)

func ParseTodoListResourceName(name string) (*TodoListResourceName, error) {
	matches := todoListRegex.FindStringSubmatch(name)
	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid todo list resource name: %q, expected: users/{user_id}/todo-lists/{list_id}", name)
	}

	uID, _ := strconv.ParseInt(matches[1], 10, 64)
	lID, _ := strconv.ParseInt(matches[2], 10, 64)

	return &TodoListResourceName{UserID: uID, TodoListID: lID}, nil
}

// users/{user_id}/todo-lists/{list_id}/todos/{todo_id}
type TodoResourceName struct {
	UserID     int64
	TodoListID int64
	TodoID     int64
}

func (n TodoResourceName) String() string {
	return fmt.Sprintf("users/%d/todo-lists/%d/todos/%d", n.UserID, n.TodoListID, n.TodoID)
}

var todoRegex = regexp.MustCompile(`^users/(\d+)/todo-lists/(\d+)/todos/(\d+)$`)

func ParseTodoResourceName(name string) (*TodoResourceName, error) {
	matches := todoRegex.FindStringSubmatch(name)
	if len(matches) != 4 {
		return nil, fmt.Errorf("invalid todo resource name: %q, expected: users/{user_id}/todo-lists/{list_id}/todos/{todo_id}", name)
	}

	uID, _ := strconv.ParseInt(matches[1], 10, 64)
	lID, _ := strconv.ParseInt(matches[2], 10, 64)
	tID, _ := strconv.ParseInt(matches[3], 10, 64)

	return &TodoResourceName{UserID: uID, TodoListID: lID, TodoID: tID}, nil
}
