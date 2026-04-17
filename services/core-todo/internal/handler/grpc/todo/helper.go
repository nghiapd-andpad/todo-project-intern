package todo

import (
	"fmt"
	"strconv"
	"strings"
)

// Helper function to parse Todo ID from Resource Name
func parseTodoID(name string) (int64, error) {
	// Chuẩn: users/{u_id}/todo-lists/{l_id}/todos/{t_id}
	parts := strings.Split(name, "/")
	if len(parts) != 6 || parts[4] != "todos" {
		return 0, fmt.Errorf("invalid resource name: %s. Expected format: users/{u_id}/todo-lists/{l_id}/todos/{t_id}", name)
	}

	id, err := strconv.ParseInt(parts[5], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid ID in resource name: %v", err)
	}
	return id, nil
}
