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

func TestCreateTodoList(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		ctx      context.Context
		stub     *stubTodoListCreator
		req      *todov1.CreateTodoListRequest
		wantCode codes.Code
		validate func(t *testing.T, stub *stubTodoListCreator, resp *todov1.CreateTodoListResponse)
	}{
		"success": {
			ctx: authCtx("1"),
			stub: &stubTodoListCreator{
				resp: &output.TodoListCreator{
					TodoList: sampleTodoList,
				},
			},
			req: &todov1.CreateTodoListRequest{
				DisplayName: "Task 1",
			},
			wantCode: codes.OK,
			validate: func(t *testing.T, stub *stubTodoListCreator, resp *todov1.CreateTodoListResponse) {
				assert.Equal(t, "Task 1", stub.gotInput.Name)
				assert.Equal(t, entity.UserID(1), stub.gotInput.OwnerID)

				assert.Equal(t, "Task 1", resp.TodoList.DisplayName)
			},
		},

		"error: missing auth": {
			ctx:      context.Background(),
			stub:     &stubTodoListCreator{},
			req:      &todov1.CreateTodoListRequest{},
			wantCode: codes.Unauthenticated,
		},

		"error: invalid auth user id": {
			ctx:      authCtx("abc"),
			stub:     &stubTodoListCreator{},
			req:      &todov1.CreateTodoListRequest{},
			wantCode: codes.Unauthenticated,
		},

		"error: usecase not found": {
			ctx: authCtx("1"),
			stub: &stubTodoListCreator{
				err: entity.NewNotFound("todo list not found"),
			},
			req: &todov1.CreateTodoListRequest{
				DisplayName: "Task 1",
			},
			wantCode: codes.NotFound,
		},

		"error: usecase internal": {
			ctx: authCtx("1"),
			stub: &stubTodoListCreator{
				err: assert.AnError,
			},
			req: &todov1.CreateTodoListRequest{
				DisplayName: "Task 1",
			},
			wantCode: codes.Internal,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			h := (&handlerBuilder{
				todoListCreator: tt.stub,
			}).build()

			resp, err := h.CreateTodoList(tt.ctx, tt.req)

			assert.Equal(t, tt.wantCode, status.Code(err))

			if tt.validate != nil {
				tt.validate(t, tt.stub, resp)
			}
		})
	}
}
