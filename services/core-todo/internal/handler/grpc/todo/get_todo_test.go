package todo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

func TestGetTodo(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		stub     *stubTodoGetter
		req      *todov1.GetTodoRequest
		wantCode codes.Code
		validate func(t *testing.T, stub *stubTodoGetter, resp *todov1.GetTodoResponse)
	}{
		"success": {
			stub: &stubTodoGetter{
				resp: &output.TodoGetter{
					Todo: sampleTodo,
				},
			},
			req: &todov1.GetTodoRequest{
				Name: "users/1/todo-lists/2/todos/3",
			},
			wantCode: codes.OK,
			validate: func(t *testing.T, stub *stubTodoGetter, resp *todov1.GetTodoResponse) {
				assert.Equal(t, entity.TodoID(3), stub.gotInput.ID)
				assert.Equal(t, "Sub Task 1", resp.Todo.Title)
			},
		},

		"error: invalid name": {
			stub:     &stubTodoGetter{},
			req:      &todov1.GetTodoRequest{Name: "bad"},
			wantCode: codes.InvalidArgument,
		},

		"error: not found": {
			stub: &stubTodoGetter{
				err: entity.NewNotFound("todo not found"),
			},
			req: &todov1.GetTodoRequest{
				Name: "users/1/todo-lists/2/todos/999",
			},
			wantCode: codes.NotFound,
		},

		"error: internal": {
			stub: &stubTodoGetter{
				err: assert.AnError,
			},
			req: &todov1.GetTodoRequest{
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
				todoGetter: tt.stub,
			}).build()

			resp, err := h.GetTodo(context.Background(), tt.req)

			assert.Equal(t, tt.wantCode, status.Code(err))

			if tt.validate != nil {
				tt.validate(t, tt.stub, resp)
			}
		})
	}
}
