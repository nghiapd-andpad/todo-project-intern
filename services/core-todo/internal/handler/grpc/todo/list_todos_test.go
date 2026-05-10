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

func TestListTodos(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		stub     *stubTodoLister
		req      *todov1.ListTodosRequest
		wantCode codes.Code
		validate func(t *testing.T, stub *stubTodoLister, resp *todov1.ListTodosResponse)
	}{
		"success": {
			stub: &stubTodoLister{
				resp: &output.TodoLister{
					Todos: []*entity.Todo{
						sampleTodo,
					},
					Total: 1,
				},
			},
			req: &todov1.ListTodosRequest{
				Parent:   "users/1/todo-lists/2",
				PageSize: 20,
			},
			wantCode: codes.OK,
			validate: func(t *testing.T, stub *stubTodoLister, resp *todov1.ListTodosResponse) {
				assert.Equal(t, int64(1), resp.Total)
				assert.Len(t, resp.Todos, 1)

				assert.Equal(t, entity.TodoListID(2), *stub.gotInput.Opts.TodoListID)
			},
		},

		"success: with filters": {
			stub: &stubTodoLister{
				resp: &output.TodoLister{},
			},
			req: &todov1.ListTodosRequest{
				Parent:         "users/1/todo-lists/2",
				TitleSearch:    "task",
				StatusFilter:   todov1.TodoStatus_TODO_STATUS_PENDING,
				PriorityFilter: todov1.Priority_PRIORITY_HIGH,
			},
			wantCode: codes.OK,
			validate: func(t *testing.T, stub *stubTodoLister, _ *todov1.ListTodosResponse) {
				assert.Equal(t, "task", *stub.gotInput.Opts.TitleSearch)
				assert.Equal(t, entity.TodoStatusPending, *stub.gotInput.Opts.Status)
				assert.Equal(t, entity.PriorityHigh, *stub.gotInput.Opts.Priority)
			},
		},

		"error: invalid parent": {
			stub:     &stubTodoLister{},
			req:      &todov1.ListTodosRequest{Parent: "bad"},
			wantCode: codes.InvalidArgument,
		},

		"error: internal": {
			stub: &stubTodoLister{
				err: assert.AnError,
			},
			req: &todov1.ListTodosRequest{
				Parent: "users/1/todo-lists/2",
			},
			wantCode: codes.Internal,
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			h := (&handlerBuilder{
				todoLister: tt.stub,
			}).build()

			resp, err := h.ListTodos(context.Background(), tt.req)

			assert.Equal(t, tt.wantCode, status.Code(err))

			if tt.validate != nil {
				tt.validate(t, tt.stub, resp)
			}
		})
	}
}
