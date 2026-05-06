package persistence_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/testutil"
)

func TestTodoListQueriesGateway_Get(t *testing.T) {
	t.Parallel()

	t.Run("success: found", func(t *testing.T) {
		t.Parallel()

		// load config
		cfg, err := config.New()
		require.NoError(t, err)

		db := testutil.NewTestDB(t, cfg)
		queryRepo := persistence.NewTodoListQueriesGateway(db)

		created := testutil.CreateTodoList(t, db, "Work Tasks", entity.UserID(1))

		got, err := queryRepo.Get(context.Background(), created.ID)

		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, created.ID, got.ID)
		assert.Equal(t, "Work Tasks", got.Name)
		assert.Equal(t, entity.UserID(1), got.OwnerID)
	})

	t.Run("not found: returns nil nil", func(t *testing.T) {
		t.Parallel()

		// load config
		cfg, err := config.New()
		require.NoError(t, err)

		db := testutil.NewTestDB(t, cfg)
		queryRepo := persistence.NewTodoListQueriesGateway(db)

		got, err := queryRepo.Get(context.Background(), entity.TodoListID(9999))

		assert.NoError(t, err)
		assert.Nil(t, got)
	})
}

func TestTodoListQueriesGateway_List(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup     func(t *testing.T, db *gorm.DB)
		opts      gateway.ListTodoListsOptions
		wantLen   int
		wantTotal int64
	}{
		"success: list by owner": {
			setup: func(t *testing.T, db *gorm.DB) {
				testutil.CreateTodoList(t, db, "Work", entity.UserID(1))
				testutil.CreateTodoList(t, db, "Personal", entity.UserID(1))
				testutil.CreateTodoList(t, db, "Other Owner", entity.UserID(2))
			},
			opts: gateway.ListTodoListsOptions{
				OwnerID: userIDPtr(entity.UserID(1)),
				Limit:   10,
			},
			wantLen:   2,
			wantTotal: 2,
		},
		"success: name search": {
			setup: func(t *testing.T, db *gorm.DB) {
				testutil.CreateTodoList(t, db, "Work Tasks", entity.UserID(1))
				testutil.CreateTodoList(t, db, "Work Projects", entity.UserID(1))
				testutil.CreateTodoList(t, db, "Personal", entity.UserID(1))
			},
			opts: gateway.ListTodoListsOptions{
				OwnerID:    userIDPtr(entity.UserID(1)),
				NameSearch: strPtr("Work"),
				Limit:      10,
			},
			wantLen:   2,
			wantTotal: 2,
		},
		"success: insensitive search": {
			setup: func(t *testing.T, db *gorm.DB) {
				testutil.CreateTodoList(t, db, "WORK TASKS", entity.UserID(1))
				testutil.CreateTodoList(t, db, "work tasks", entity.UserID(1))
			},
			opts: gateway.ListTodoListsOptions{
				OwnerID:    userIDPtr(entity.UserID(1)),
				NameSearch: strPtr("work"),
				Limit:      10,
			},
			wantLen:   2,
			wantTotal: 2,
		},
		"success: pagination": {
			setup: func(t *testing.T, db *gorm.DB) {
				for i := 0; i < 5; i++ {
					testutil.CreateTodoList(t, db,
						fmt.Sprintf("List %d", i),
						entity.UserID(1),
					)
				}
			},
			opts: gateway.ListTodoListsOptions{
				OwnerID: userIDPtr(entity.UserID(1)),
				Limit:   2,
				Offset:  2,
			},
			wantLen:   2,
			wantTotal: 5,
		},
		"success: empty": {
			setup: func(t *testing.T, db *gorm.DB) {},
			opts: gateway.ListTodoListsOptions{
				OwnerID: userIDPtr(entity.UserID(1)),
				Limit:   10,
			},
			wantLen:   0,
			wantTotal: 0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// load config
			cfg, err := config.New()
			require.NoError(t, err)

			db := testutil.NewTestDB(t, cfg)
			tt.setup(t, db)
			queryRepo := persistence.NewTodoListQueriesGateway(db)

			got, total, err := queryRepo.List(context.Background(), tt.opts)

			require.NoError(t, err)
			assert.Len(t, got, tt.wantLen)
			assert.Equal(t, tt.wantTotal, total)
		})
	}
}
