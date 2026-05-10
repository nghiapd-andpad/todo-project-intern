package todo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

func TestUpdateTodoList(t *testing.T) {
	t.Parallel()

	newName := "New List"

	tests := map[string]struct {
		stub     *stubTodoListUpdater
		req      *todov1.UpdateTodoListRequest
		wantCode codes.Code
		validate func(t *testing.T, stub *stubTodoListUpdater)
	}{
		"success": {
			stub: &stubTodoListUpdater{
				resp: &output.TodoListUpdater{
					TodoList: sampleTodoList,
				},
			},
			req: &todov1.UpdateTodoListRequest{
				TodoList: &todov1.TodoList{
					Name:        "users/1/todo-lists/2",
					DisplayName: newName,
				},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"display_name"},
				},
			},
			wantCode: codes.OK,
			validate: func(t *testing.T, stub *stubTodoListUpdater) {
				assert.Equal(t, entity.TodoListID(2), stub.gotInput.ID)
				assert.Equal(t, newName, *stub.gotInput.Name)
			},
		},

		"error: missing todo_list": {
			stub: &stubTodoListUpdater{},
			req: &todov1.UpdateTodoListRequest{
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"display_name"},
				},
			},
			wantCode: codes.InvalidArgument,
		},

		"error: missing update_mask": {
			stub: &stubTodoListUpdater{},
			req: &todov1.UpdateTodoListRequest{
				TodoList: &todov1.TodoList{
					Name: "users/1/todo-lists/2",
				},
			},
			wantCode: codes.InvalidArgument,
		},

		"error: invalid name": {
			stub: &stubTodoListUpdater{},
			req: &todov1.UpdateTodoListRequest{
				TodoList: &todov1.TodoList{
					Name: "bad",
				},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"display_name"},
				},
			},
			wantCode: codes.InvalidArgument,
		},

		"error: not found": {
			stub: &stubTodoListUpdater{
				err: entity.NewNotFound("todo list not found"),
			},
			req: &todov1.UpdateTodoListRequest{
				TodoList: &todov1.TodoList{
					Name: "users/1/todo-lists/999",
				},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"display_name"},
				},
			},
			wantCode: codes.NotFound,
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			h := (&handlerBuilder{
				todoListUpdater: tt.stub,
			}).build()

			_, err := h.UpdateTodoList(context.Background(), tt.req)

			assert.Equal(t, tt.wantCode, status.Code(err))

			if tt.validate != nil {
				tt.validate(t, tt.stub)
			}
		})
	}
}
