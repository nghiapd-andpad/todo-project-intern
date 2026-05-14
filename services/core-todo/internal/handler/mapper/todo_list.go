package mapper

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/resourcename"
	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
)

func TodoListToPb(ent *entity.TodoList) *todov1.TodoList {
	if ent == nil {
		return nil
	}

	return &todov1.TodoList{
		Name: resourcename.TodoListResourceName{
			UserID:     int64(ent.OwnerID),
			TodoListID: int64(ent.ID),
		}.String(),
		DisplayName: ent.Name,
		CreatedAt:   timestamppb.New(ent.CreatedAt),
		UpdatedAt:   timestamppb.New(ent.UpdatedAt),
	}
}
