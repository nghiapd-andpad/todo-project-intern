package resourcename

import (
	"fmt"
	"regexp"
	"strconv"
)

type TodoResourceName struct {
	UserID     int64
	TodoListID int64
	TodoID     int64
}

var todoRegex = regexp.MustCompile(`^users/(\d+)/todo-lists/(\d+)/todos/(\d+)$`)

func ParseTodoResourceName(name string) (*TodoResourceName, error) {
	matches := todoRegex.FindStringSubmatch(name)
	if len(matches) != 4 {
		return nil, fmt.Errorf("invalid resource name: %s", name)
	}

	uID, _ := strconv.ParseInt(matches[1], 10, 64)
	lID, _ := strconv.ParseInt(matches[2], 10, 64)
	tID, _ := strconv.ParseInt(matches[3], 10, 64)

	return &TodoResourceName{UserID: uID, TodoListID: lID, TodoID: tID}, nil
}
