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

func TestCreateTodo(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		ctx      context.Context
		stub     *stubTodoCreator
		req      *todov1.CreateTodoRequest
		wantCode codes.Code
		validate func(t *testing.T, stub *stubTodoCreator, resp *todov1.CreateTodoResponse)
	}{
		"success": {
			ctx: authCtx("1"),
			stub: &stubTodoCreator{
				resp: &output.TodoCreator{
					Todo: sampleTodo,
				},
			},
			req: &todov1.CreateTodoRequest{
				Parent: "users/1/todo-lists/2",
				Title:  "Sub Task 1",
			},
			wantCode: codes.OK,
			validate: func(t *testing.T, stub *stubTodoCreator, resp *todov1.CreateTodoResponse) {
				assert.Equal(t, entity.TodoListID(2), stub.gotInput.TodoListID)
				assert.Equal(t, entity.UserID(1), stub.gotInput.CreatorID)
				assert.Equal(t, "Sub Task 1", stub.gotInput.Title)
				assert.Equal(t, "Sub Task 1", resp.Todo.Title)
			},
		},

		"success: optional fields": {
			ctx: authCtx("1"),
			stub: &stubTodoCreator{
				resp: &output.TodoCreator{
					Todo: sampleTodo,
				},
			},
			req: &todov1.CreateTodoRequest{
				Parent:      "users/1/todo-lists/2",
				Title:       "Sub Task 1",
				Description: "Need money",
				DueDate:     "2026-05-10",
				AssigneeId:  99,
				Priority:    todov1.Priority_PRIORITY_HIGH,
			},
			wantCode: codes.OK,
			validate: func(t *testing.T, stub *stubTodoCreator, _ *todov1.CreateTodoResponse) {
				assert.Equal(t, "Need money", *stub.gotInput.Description)
				assert.Equal(t, "2026-05-10", *stub.gotInput.DueDate)
				assert.Equal(t, entity.UserID(99), *stub.gotInput.AssigneeID)
				assert.Equal(t, entity.PriorityHigh, stub.gotInput.Priority)
			},
		},

		"error: invalid parent": {
			ctx:      authCtx("1"),
			stub:     &stubTodoCreator{},
			req:      &todov1.CreateTodoRequest{Parent: "bad"},
			wantCode: codes.InvalidArgument,
			validate: func(t *testing.T, stub *stubTodoCreator, _ *todov1.CreateTodoResponse) {
				assert.Nil(t, stub.gotInput)
			},
		},

		"error: missing auth": {
			ctx:      context.Background(),
			stub:     &stubTodoCreator{},
			req:      &todov1.CreateTodoRequest{Parent: "users/1/todo-lists/2"},
			wantCode: codes.Unauthenticated,
			validate: func(t *testing.T, stub *stubTodoCreator, _ *todov1.CreateTodoResponse) {
				assert.Nil(t, stub.gotInput)
			},
		},

		"error: invalid auth user id": {
			ctx:      authCtx("abc"),
			stub:     &stubTodoCreator{},
			req:      &todov1.CreateTodoRequest{Parent: "users/1/todo-lists/2"},
			wantCode: codes.Unauthenticated,
		},

		"error: usecase not found": {
			ctx: authCtx("1"),
			stub: &stubTodoCreator{
				err: entity.NewNotFound("todo list not found"),
			},
			req: &todov1.CreateTodoRequest{
				Parent: "users/1/todo-lists/2",
			},
			wantCode: codes.NotFound,
		},

		"error: usecase internal": {
			ctx: authCtx("1"),
			stub: &stubTodoCreator{
				err: assert.AnError,
			},
			req: &todov1.CreateTodoRequest{
				Parent: "users/1/todo-lists/2",
			},
			wantCode: codes.Internal,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			h := (&handlerBuilder{
				todoCreator: tt.stub,
			}).build()

			resp, err := h.CreateTodo(tt.ctx, tt.req)

			assert.Equal(t, tt.wantCode, status.Code(err))

			if tt.validate != nil {
				tt.validate(t, tt.stub, resp)
			}
		})
	}
}
