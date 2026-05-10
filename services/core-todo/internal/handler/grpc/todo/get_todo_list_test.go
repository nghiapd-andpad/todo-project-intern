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

func TestGetTodoList(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		stub     *stubTodoListGetter
		req      *todov1.GetTodoListRequest
		wantCode codes.Code
		validate func(t *testing.T, stub *stubTodoListGetter, resp *todov1.GetTodoListResponse)
	}{
		"success": {
			stub: &stubTodoListGetter{
				resp: &output.TodoListGetter{
					TodoList: sampleTodoList,
				},
			},
			req: &todov1.GetTodoListRequest{
				Name: "users/1/todo-lists/2",
			},
			wantCode: codes.OK,
			validate: func(t *testing.T, stub *stubTodoListGetter, resp *todov1.GetTodoListResponse) {
				assert.Equal(t, entity.TodoListID(2), stub.gotInput.ID)
				assert.Equal(t, "Task 1", resp.TodoList.DisplayName)
			},
		},

		"error: invalid name": {
			stub:     &stubTodoListGetter{},
			req:      &todov1.GetTodoListRequest{Name: "bad"},
			wantCode: codes.InvalidArgument,
			validate: func(t *testing.T, stub *stubTodoListGetter, _ *todov1.GetTodoListResponse) {
				assert.Nil(t, stub.gotInput)
			},
		},

		"error: not found": {
			stub: &stubTodoListGetter{
				err: entity.NewNotFound("todo list not found"),
			},
			req: &todov1.GetTodoListRequest{
				Name: "users/1/todo-lists/999",
			},
			wantCode: codes.NotFound,
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			h := (&handlerBuilder{
				todoListGetter: tt.stub,
			}).build()

			resp, err := h.GetTodoList(context.Background(), tt.req)

			assert.Equal(t, tt.wantCode, status.Code(err))

			if tt.wantCode == codes.OK {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}

			if tt.validate != nil {
				tt.validate(t, tt.stub, resp)
			}
		})
	}
}
