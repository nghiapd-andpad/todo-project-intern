package todo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
)

func TestDeleteTodo(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		stub     *stubTodoDeleter
		req      *todov1.DeleteTodoRequest
		wantCode codes.Code
		validate func(t *testing.T, stub *stubTodoDeleter)
	}{
		"success": {
			stub: &stubTodoDeleter{},
			req: &todov1.DeleteTodoRequest{
				Name: "users/1/todo-lists/2/todos/3",
			},
			wantCode: codes.OK,
			validate: func(t *testing.T, stub *stubTodoDeleter) {
				assert.Equal(t, entity.TodoID(3), stub.gotInput.ID)
			},
		},

		"error: invalid name": {
			stub:     &stubTodoDeleter{},
			req:      &todov1.DeleteTodoRequest{Name: "bad"},
			wantCode: codes.InvalidArgument,
		},

		"error: not found": {
			stub: &stubTodoDeleter{
				err: entity.NewNotFound("todo not found"),
			},
			req: &todov1.DeleteTodoRequest{
				Name: "users/1/todo-lists/2/todos/999",
			},
			wantCode: codes.NotFound,
		},

		"error: internal": {
			stub: &stubTodoDeleter{
				err: assert.AnError,
			},
			req: &todov1.DeleteTodoRequest{
				Name: "users/1/todo-lists/2/todos/3",
			},
			wantCode: codes.Internal,
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			h := (&handlerBuilder{
				todoDeleter: tt.stub,
			}).build()

			resp, err := h.DeleteTodo(context.Background(), tt.req)

			assert.Equal(t, tt.wantCode, status.Code(err))

			if tt.wantCode == codes.OK {
				assert.NotNil(t, resp)
			}

			if tt.validate != nil {
				tt.validate(t, tt.stub)
			}
		})
	}
}
