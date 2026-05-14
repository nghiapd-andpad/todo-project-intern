package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/mock"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/service"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/testutil"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/output"
)

func TestTodoCreator_Create(t *testing.T) {
	t.Parallel()

	var (
		ctx         = context.Background()
		now         = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		todoListID  = entity.TodoListID(2)
		requesterID = entity.UserID(1)
		ownerID     = entity.UserID(1)
		assigneeID  = entity.UserID(3)
		dueDate     = time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)

		existingTodoList = &entity.TodoList{
			ID:      todoListID,
			Name:    "Work Tasks",
			OwnerID: ownerID,
		}

		createdEntity = &entity.Todo{
			ID:         entity.TodoID(10),
			TodoListID: todoListID,
			Title:      "Unit Test Create Todo",
			Status:     entity.TodoStatusPending,
			Priority:   entity.PriorityMedium,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		validInput = &input.TodoCreator{
			TodoListID:  todoListID,
			Title:       "Unit Test Create Todo",
			Priority:    entity.PriorityMedium,
			RequesterID: requesterID,
		}

		defaultCfg = &config.Config{
			TodoBlacklistEnabled: false,
		}

		blacklistCfg = &config.Config{
			TodoBlacklistEnabled: true,
			TodoTitleBlacklist:   []string{"spam", "troll"},
		}
	)

	type fields struct {
		mockTodoListQueries *mock.MockTodoListQueriesGateway
		mockTodoCommands    *mock.MockTodoCommandsGateway
	}

	tests := map[string]struct {
		cfg      *config.Config
		prepare  func(f *fields)
		input    *input.TodoCreator
		expected *output.TodoCreator
		wantErr  bool
		errCode  entity.ErrorCode
	}{
		"success: create with required fields": {
			cfg: defaultCfg,
			prepare: func(f *fields) {
				f.mockTodoListQueries.EXPECT().
					Get(gomock.Any(), todoListID).
					Return(existingTodoList, nil)

				f.mockTodoCommands.EXPECT().
					Create(gomock.Any(), &entity.Todo{
						TodoListID: todoListID,
						Title:      "Unit Test Create Todo",
						Status:     entity.TodoStatusPending,
						Priority:   entity.PriorityMedium,
					}).
					Return(createdEntity, nil)
			},
			input:    validInput,
			expected: &output.TodoCreator{Todo: createdEntity},
		},

		"success: create with optional fields": {
			cfg: defaultCfg,
			prepare: func(f *fields) {
				f.mockTodoListQueries.EXPECT().
					Get(gomock.Any(), todoListID).
					Return(existingTodoList, nil)

				f.mockTodoCommands.EXPECT().
					Create(gomock.Any(), &entity.Todo{
						TodoListID:  todoListID,
						Title:       "Unit Test Create Todo",
						Description: testutil.StrPtr("description"),
						Status:      entity.TodoStatusPending,
						Priority:    entity.PriorityHigh,
						DueDate:     &dueDate,
						AssigneeID:  &assigneeID,
					}).
					Return(createdEntity, nil)
			},
			input: &input.TodoCreator{
				TodoListID:  todoListID,
				Title:       "Unit Test Create Todo",
				Description: testutil.StrPtr("description"),
				Priority:    entity.PriorityHigh,
				DueDate:     &dueDate,
				AssigneeID:  &assigneeID,
				RequesterID: requesterID,
			},
			expected: &output.TodoCreator{Todo: createdEntity},
		},

		"success: blacklist enabled and title is valid": {
			cfg: blacklistCfg,
			prepare: func(f *fields) {
				f.mockTodoListQueries.EXPECT().
					Get(gomock.Any(), todoListID).
					Return(existingTodoList, nil)

				f.mockTodoCommands.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(createdEntity, nil)
			},
			input: &input.TodoCreator{
				TodoListID:  todoListID,
				Title:       "Normal task",
				Priority:    entity.PriorityMedium,
				RequesterID: requesterID,
			},
			expected: &output.TodoCreator{Todo: createdEntity},
		},

		"success: blacklist disabled and title contains blacklisted word": {
			cfg: defaultCfg,
			prepare: func(f *fields) {
				f.mockTodoListQueries.EXPECT().
					Get(gomock.Any(), todoListID).
					Return(existingTodoList, nil)

				f.mockTodoCommands.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(createdEntity, nil)
			},
			input: &input.TodoCreator{
				TodoListID:  todoListID,
				Title:       "spam task",
				Priority:    entity.PriorityMedium,
				RequesterID: requesterID,
			},
			expected: &output.TodoCreator{Todo: createdEntity},
		},

		"error: blacklist enabled and title contains blacklisted word": {
			cfg:     blacklistCfg,
			prepare: func(f *fields) {},
			input: &input.TodoCreator{
				TodoListID:  todoListID,
				Title:       "spam task",
				Priority:    entity.PriorityMedium,
				RequesterID: requesterID,
			},
			wantErr: true,
			errCode: entity.ErrInvalidParameter,
		},

		"error: blacklist is case-insensitive": {
			cfg:     blacklistCfg,
			prepare: func(f *fields) {},
			input: &input.TodoCreator{
				TodoListID:  todoListID,
				Title:       "SPAM task",
				Priority:    entity.PriorityMedium,
				RequesterID: requesterID,
			},
			wantErr: true,
			errCode: entity.ErrInvalidParameter,
		},

		"error: todo list query gateway error": {
			cfg: defaultCfg,
			prepare: func(f *fields) {
				f.mockTodoListQueries.EXPECT().
					Get(gomock.Any(), todoListID).
					Return(nil, fmt.Errorf("db connection lost"))
			},
			input:   validInput,
			wantErr: true,
		},

		"error: todo list not found": {
			cfg: defaultCfg,
			prepare: func(f *fields) {
				f.mockTodoListQueries.EXPECT().
					Get(gomock.Any(), todoListID).
					Return(nil, nil)
			},
			input:   validInput,
			wantErr: true,
			errCode: entity.ErrNotFound,
		},

		"error: requester is not todo list owner": {
			cfg: defaultCfg,
			prepare: func(f *fields) {
				f.mockTodoListQueries.EXPECT().
					Get(gomock.Any(), todoListID).
					Return(&entity.TodoList{
						ID:      todoListID,
						Name:    "Work Tasks",
						OwnerID: entity.UserID(999),
					}, nil)
			},
			input:   validInput,
			wantErr: true,
			errCode: entity.ErrAuthZ,
		},

		"error: todo command gateway error": {
			cfg: defaultCfg,
			prepare: func(f *fields) {
				f.mockTodoListQueries.EXPECT().
					Get(gomock.Any(), todoListID).
					Return(existingTodoList, nil)

				f.mockTodoCommands.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("db connection lost"))
			},
			input:   validInput,
			wantErr: true,
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			f := &fields{
				mockTodoListQueries: mock.NewMockTodoListQueriesGateway(ctrl),
				mockTodoCommands:    mock.NewMockTodoCommandsGateway(ctrl),
			}

			tt.prepare(f)

			sut := service.NewTodoCreator(
				f.mockTodoListQueries,
				f.mockTodoCommands,
				tt.cfg,
			)

			got, err := sut.Create(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)

				if tt.errCode != "" {
					var appErr *entity.AppError
					assert.ErrorAs(t, err, &appErr)
					assert.Equal(t, tt.errCode, appErr.Code)
				}

				return
			}

			assert.NoError(t, err)

			if diff := cmp.Diff(tt.expected, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
