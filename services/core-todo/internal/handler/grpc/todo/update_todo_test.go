package todo_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

func TestUpdateTodo(t *testing.T) {
	t.Parallel()

	newTitle := "New Title"

	tests := map[string]struct {
		stub     *stubTodoUpdater
		req      *todov1.UpdateTodoRequest
		wantCode codes.Code
		validate func(t *testing.T, stub *stubTodoUpdater)
	}{
		"success: update title": {
			stub: &stubTodoUpdater{
				resp: &output.TodoUpdater{
					Todo: sampleTodo,
				},
			},
			req: &todov1.UpdateTodoRequest{
				Todo: &todov1.Todo{
					Name:  "users/1/todo-lists/2/todos/3",
					Title: newTitle,
				},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"title"},
				},
			},
			wantCode: codes.OK,
			validate: func(t *testing.T, stub *stubTodoUpdater) {
				assert.Equal(t, entity.TodoID(3), stub.gotInput.ID)
				assert.Equal(t, &newTitle, stub.gotInput.Title)
			},
		},

		"success: update description": {
			stub: &stubTodoUpdater{
				resp: &output.TodoUpdater{Todo: sampleTodo},
			},
			req: &todov1.UpdateTodoRequest{
				Todo: &todov1.Todo{
					Name:        "users/1/todo-lists/2/todos/3",
					Description: "new description",
				},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"description"},
				},
			},
			wantCode: codes.OK,
			validate: func(t *testing.T, stub *stubTodoUpdater) {
				assert.Equal(t, "new description", *stub.gotInput.Description)
			},
		},

		"success: update status": {
			stub: &stubTodoUpdater{
				resp: &output.TodoUpdater{Todo: sampleTodo},
			},
			req: &todov1.UpdateTodoRequest{
				Todo: &todov1.Todo{
					Name:   "users/1/todo-lists/2/todos/3",
					Status: todov1.TodoStatus_TODO_STATUS_DONE,
				},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"status"},
				},
			},
			wantCode: codes.OK,
			validate: func(t *testing.T, stub *stubTodoUpdater) {
				assert.Equal(t, entity.TodoStatusDone, *stub.gotInput.Status)
			},
		},

		"success: update priority": {
			stub: &stubTodoUpdater{
				resp: &output.TodoUpdater{Todo: sampleTodo},
			},
			req: &todov1.UpdateTodoRequest{
				Todo: &todov1.Todo{
					Name:     "users/1/todo-lists/2/todos/3",
					Priority: todov1.Priority_PRIORITY_HIGH,
				},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"priority"},
				},
			},
			wantCode: codes.OK,
			validate: func(t *testing.T, stub *stubTodoUpdater) {
				assert.Equal(t, entity.PriorityHigh, *stub.gotInput.Priority)
			},
		},

		"success: update due_date": {
			stub: &stubTodoUpdater{
				resp: &output.TodoUpdater{Todo: sampleTodo},
			},
			req: &todov1.UpdateTodoRequest{
				Todo: &todov1.Todo{
					Name:    "users/1/todo-lists/2/todos/3",
					DueDate: timestamppb.New(time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)),
				},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"due_date"},
				},
			},
			wantCode: codes.OK,
			validate: func(t *testing.T, stub *stubTodoUpdater) {
				assert.Equal(t, "2026-05-10", *stub.gotInput.DueDate)
			},
		},

		"success: update assignee_id": {
			stub: &stubTodoUpdater{
				resp: &output.TodoUpdater{Todo: sampleTodo},
			},
			req: &todov1.UpdateTodoRequest{
				Todo: &todov1.Todo{
					Name:       "users/1/todo-lists/2/todos/3",
					AssigneeId: 55,
				},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"assignee_id"},
				},
			},
			wantCode: codes.OK,
			validate: func(t *testing.T, stub *stubTodoUpdater) {
				assert.Equal(t, entity.UserID(55), *stub.gotInput.AssigneeID)
			},
		},

		"error: missing todo": {
			stub: &stubTodoUpdater{},
			req: &todov1.UpdateTodoRequest{
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"title"},
				},
			},
			wantCode: codes.InvalidArgument,
		},

		"error: missing update_mask": {
			stub: &stubTodoUpdater{},
			req: &todov1.UpdateTodoRequest{
				Todo: &todov1.Todo{
					Name: "users/1/todo-lists/2/todos/3",
				},
			},
			wantCode: codes.InvalidArgument,
		},

		"error: invalid name": {
			stub: &stubTodoUpdater{},
			req: &todov1.UpdateTodoRequest{
				Todo: &todov1.Todo{
					Name: "bad",
				},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"title"},
				},
			},
			wantCode: codes.InvalidArgument,
		},

		"error: not found": {
			stub: &stubTodoUpdater{
				err: entity.NewNotFound("todo not found"),
			},
			req: &todov1.UpdateTodoRequest{
				Todo: &todov1.Todo{
					Name: "users/1/todo-lists/2/todos/999",
				},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"title"},
				},
			},
			wantCode: codes.NotFound,
		},

		"error: internal": {
			stub: &stubTodoUpdater{
				err: assert.AnError,
			},
			req: &todov1.UpdateTodoRequest{
				Todo: &todov1.Todo{
					Name: "users/1/todo-lists/2/todos/3",
				},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"title"},
				},
			},
			wantCode: codes.Internal,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			h := (&handlerBuilder{
				todoUpdater: tt.stub,
			}).build()

			_, err := h.UpdateTodo(context.Background(), tt.req)

			assert.Equal(t, tt.wantCode, status.Code(err))

			if tt.validate != nil {
				tt.validate(t, tt.stub)
			}
		})
	}
}
