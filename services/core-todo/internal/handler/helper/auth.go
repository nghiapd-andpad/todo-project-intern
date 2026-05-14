package helper

import (
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
)

func ExtractRequesterID(userIDStr string, ok bool) (entity.UserID, error) {
	if !ok || userIDStr == "" {
		return 0, status.Error(codes.Unauthenticated, "missing user id in context")
	}
	id, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, status.Error(codes.Unauthenticated, "invalid user id in context")
	}
	return entity.UserID(id), nil
}
