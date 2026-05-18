package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/mock"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/service"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/output"
)

func TestTodoListCreator_Create(t *testing.T) {
	t.Parallel()

	var (
		now         = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		requesterID = entity.UserID(1)

		validInput = &input.TodoListCreator{
			Name:        "Task 1",
			RequesterID: requesterID,
		}

		createdEntity = &entity.TodoList{
			ID:        entity.TodoListID(1),
			Name:      "Task 1",
			OwnerID:   requesterID,
			CreatedAt: now,
			UpdatedAt: now,
		}
	)

	type fields struct {
		mockCommands *mock.MockTodoListCommandsGateway
	}

	tests := map[string]struct {
		prepare  func(f *fields, ctx context.Context)
		input    *input.TodoListCreator
		expected *output.TodoListCreator
		wantErr  bool
	}{
		"success: create todo list": {
			prepare: func(f *fields, ctx context.Context) {
				f.mockCommands.EXPECT().
					Create(ctx, &entity.TodoList{
						Name:    "Task 1",
						OwnerID: requesterID,
					}).
					Return(createdEntity, nil)
			},
			input:    validInput,
			expected: &output.TodoListCreator{TodoList: createdEntity},
		},

		"success: create with empty name": {
			prepare: func(f *fields, ctx context.Context) {
				f.mockCommands.EXPECT().
					Create(ctx, &entity.TodoList{
						Name:    "",
						OwnerID: requesterID,
					}).
					Return(&entity.TodoList{ID: 2, Name: "", OwnerID: requesterID}, nil)
			},
			input: &input.TodoListCreator{
				Name:        "",
				RequesterID: requesterID,
			},
			expected: &output.TodoListCreator{
				TodoList: &entity.TodoList{ID: 2, Name: "", OwnerID: requesterID},
			},
		},

		"error: command gateway error": {
			prepare: func(f *fields, ctx context.Context) {
				f.mockCommands.EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil, fmt.Errorf("db connection lost"))
			},
			input:   validInput,
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			ctrl := gomock.NewController(t)

			f := &fields{
				mockCommands: mock.NewMockTodoListCommandsGateway(ctrl),
			}

			tt.prepare(f, ctx)

			sut := service.NewTodoListCreator(f.mockCommands)

			got, err := sut.Create(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if diff := cmp.Diff(tt.expected, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
