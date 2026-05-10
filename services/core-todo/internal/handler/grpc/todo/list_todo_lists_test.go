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

func TestListTodoLists(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		stub     *stubTodoListLister
		req      *todov1.ListTodoListsRequest
		wantCode codes.Code
		validate func(t *testing.T, stub *stubTodoListLister, resp *todov1.ListTodoListsResponse)
	}{
		"success": {
			stub: &stubTodoListLister{
				resp: &output.TodoListLister{
					TodoLists: []*entity.TodoList{
						sampleTodoList,
					},
					Total: 1,
				},
			},
			req: &todov1.ListTodoListsRequest{
				Parent:   "users/1",
				PageSize: 20,
			},
			wantCode: codes.OK,
			validate: func(t *testing.T, _ *stubTodoListLister, resp *todov1.ListTodoListsResponse) {
				assert.Len(t, resp.TodoLists, 1)
				assert.Equal(t, int64(1), resp.Total)
			},
		},

		"error: invalid parent": {
			stub:     &stubTodoListLister{},
			req:      &todov1.ListTodoListsRequest{Parent: "bad"},
			wantCode: codes.InvalidArgument,
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			h := (&handlerBuilder{
				todoListLister: tt.stub,
			}).build()

			resp, err := h.ListTodoLists(context.Background(), tt.req)

			assert.Equal(t, tt.wantCode, status.Code(err))

			if tt.wantCode == codes.OK {
				assert.NotNil(t, resp)
			}

			if tt.validate != nil {
				tt.validate(t, tt.stub, resp)
			}
		})
	}
}
