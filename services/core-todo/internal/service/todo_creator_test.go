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
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/output"
)

func TestTodoCreator_Create(t *testing.T) {
	t.Parallel()

	var (
		ctx        = context.Background()
		now        = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		todoListID = entity.TodoListID(2)
		creatorID  = entity.UserID(1)
		validDue   = "2026-05-01"
		parsedDue  = time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)

		createdEntity = &entity.Todo{
			ID:         entity.TodoID(10),
			TodoListID: todoListID,
			Title:      "Unit Test Create Todo",
			Status:     entity.TodoStatusPending,
			Priority:   entity.PriorityMedium,
			CreatorID:  creatorID,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		validInput = &input.TodoCreator{
			TodoListID: todoListID,
			Title:      "Unit Test Create Todo",
			Priority:   entity.PriorityMedium,
			CreatorID:  creatorID,
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
		mockCommands *mock.MockTodoCommandsGateway
	}

	tests := map[string]struct {
		cfg      *config.Config
		prepare  func(f *fields)
		input    *input.TodoCreator
		expected *output.TodoCreator
		wantErr  bool
		errCode  entity.ErrorCode
	}{
		// Happy path — verify required fields
		"success: create with required fields": {
			cfg: defaultCfg,
			prepare: func(f *fields) {
				f.mockCommands.EXPECT().
					Create(gomock.Any(), &entity.Todo{
						TodoListID: todoListID,
						Title:      "Unit Test Create Todo",
						Status:     entity.TodoStatusPending,
						Priority:   entity.PriorityMedium,
						CreatorID:  creatorID,
					}).
					Return(createdEntity, nil)
			},
			wantErr:  false,
			input:    validInput,
			expected: &output.TodoCreator{Todo: createdEntity},
		},

		// Happy path — parse DueDate
		"success: create with due_date": {
			cfg: defaultCfg,
			prepare: func(f *fields) {
				f.mockCommands.EXPECT().
					Create(gomock.Any(), &entity.Todo{
						TodoListID: todoListID,
						Title:      "Unit Test Create Todo",
						Status:     entity.TodoStatusPending,
						Priority:   entity.PriorityMedium,
						CreatorID:  creatorID,
						DueDate:    &parsedDue,
					}).
					Return(createdEntity, nil)
			},
			input: &input.TodoCreator{
				TodoListID: todoListID,
				Title:      "Unit Test Create Todo",
				Priority:   entity.PriorityMedium,
				CreatorID:  creatorID,
				DueDate:    &validDue,
			},
			wantErr:  false,
			expected: &output.TodoCreator{Todo: createdEntity},
		},

		// Error path — DueDate wrong format
		"error: invalid due_date format": {
			cfg: defaultCfg,
			prepare: func(f *fields) {
			},
			input: &input.TodoCreator{
				TodoListID: todoListID,
				Title:      "Unit Test Create Todo",
				CreatorID:  creatorID,
				DueDate:    func() *string { s := "01/05/2026"; return &s }(),
			},
			wantErr: true,
			errCode: entity.ErrInvalidParameter,
		},

		// Error path — gateway DB error
		"error: gateway db error": {
			cfg: defaultCfg,
			prepare: func(f *fields) {
				f.mockCommands.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("db connection lost"))
			},
			input:   validInput,
			wantErr: true,
		},

		// Feature flag ON
		"success: flag ON, title valid and not blocked": {
			cfg: blacklistCfg,
			prepare: func(f *fields) {
				f.mockCommands.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(createdEntity, nil)
			},
			input: &input.TodoCreator{
				TodoListID: todoListID,
				Title:      "Task 1",
				Priority:   entity.PriorityMedium,
				CreatorID:  creatorID,
			},
			expected: &output.TodoCreator{Todo: createdEntity},
		},

		"error: flag ON, title contains blacklisted words": {
			cfg:     blacklistCfg,
			prepare: func(f *fields) {},
			input: &input.TodoCreator{
				TodoListID: todoListID,
				Title:      "spam ahihi",
				Priority:   entity.PriorityMedium,
				CreatorID:  creatorID,
			},
			wantErr: true,
			errCode: entity.ErrInvalidParameter,
		},

		"error: flag ON, title bypass in uppercase still gets blocked.": {
			cfg:     blacklistCfg,
			prepare: func(f *fields) {},
			input: &input.TodoCreator{
				TodoListID: todoListID,
				Title:      "SPAm Ahihi",
				Priority:   entity.PriorityMedium,
				CreatorID:  creatorID,
			},
			wantErr: true,
			errCode: entity.ErrInvalidParameter,
		},

		"success: flag OFF, titles containing prohibited words can still be created.": {
			cfg: defaultCfg,
			prepare: func(f *fields) {
				f.mockCommands.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(createdEntity, nil)
			},
			input: &input.TodoCreator{
				TodoListID: todoListID,
				Title:      "spam ahihi",
				Priority:   entity.PriorityMedium,
				CreatorID:  creatorID,
			},
			expected: &output.TodoCreator{Todo: createdEntity},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			f := &fields{mockCommands: mock.NewMockTodoCommandsGateway(ctrl)}
			tt.prepare(f)

			sut := service.NewTodoCreator(f.mockCommands, tt.cfg)
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
