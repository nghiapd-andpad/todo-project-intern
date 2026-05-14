package persistence_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/testutil"
)

func TestTodoListQueriesGateway_Get(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup func(t *testing.T) (*gorm.DB, *entity.TodoList)
		test  func(
			t *testing.T,
			repo *persistence.TodoListQueriesGateway,
			existingTodoList *entity.TodoList,
		)
	}{
		"success: found": {
			setup: func(t *testing.T) (*gorm.DB, *entity.TodoList) {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				existingTodoList := testutil.CreateTodoList(
					t,
					db,
					"Work Tasks",
					entity.UserID(1),
				)

				return db, existingTodoList
			},
			test: func(
				t *testing.T,
				repo *persistence.TodoListQueriesGateway,
				existingTodoList *entity.TodoList,
			) {
				got, err := repo.Get(
					context.Background(),
					existingTodoList.ID,
				)

				require.NoError(t, err)
				require.NotNil(t, got)

				assert.Equal(t, existingTodoList.ID, got.ID)
				assert.Equal(t, "Work Tasks", got.Name)
				assert.Equal(t, entity.UserID(1), got.OwnerID)
			},
		},

		"not found: returns nil nil": {
			setup: func(t *testing.T) (*gorm.DB, *entity.TodoList) {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				return db, nil
			},
			test: func(
				t *testing.T,
				repo *persistence.TodoListQueriesGateway,
				_ *entity.TodoList,
			) {
				got, err := repo.Get(
					context.Background(),
					entity.TodoListID(9999),
				)

				assert.NoError(t, err)
				assert.Nil(t, got)
			},
		},

		"not found: soft-deleted todo list returns nil nil": {
			setup: func(t *testing.T) (*gorm.DB, *entity.TodoList) {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				existingTodoList := testutil.CreateTodoList(
					t,
					db,
					"Deleted List",
					entity.UserID(1),
				)

				cmdRepo := persistence.NewTodoListCommandsGateway(db)

				err := cmdRepo.Delete(context.Background(), existingTodoList.ID)
				require.NoError(t, err)

				return db, existingTodoList
			},
			test: func(
				t *testing.T,
				repo *persistence.TodoListQueriesGateway,
				existingTodoList *entity.TodoList,
			) {
				got, err := repo.Get(
					context.Background(),
					existingTodoList.ID,
				)

				require.NoError(t, err)
				assert.Nil(t, got)
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, existingTodoList := tt.setup(t)

			repo := persistence.NewTodoListQueriesGateway(db)

			tt.test(t, repo, existingTodoList)
		})
	}
}

func TestTodoListQueriesGateway_List(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup     func(t *testing.T) *gorm.DB
		opts      gatewayinput.ListTodoListsOptions
		wantLen   int
		wantTotal int64
	}{
		"success: list by owner": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				testutil.CreateTodoList(t, db, "Work", entity.UserID(1))
				testutil.CreateTodoList(t, db, "Personal", entity.UserID(1))
				testutil.CreateTodoList(t, db, "Other Owner", entity.UserID(2))

				return db
			},
			opts: gatewayinput.ListTodoListsOptions{
				OwnerID: testutil.UserIDPtr(entity.UserID(1)),
				Limit:   10,
			},
			wantLen:   2,
			wantTotal: 2,
		},

		"success: name search": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				testutil.CreateTodoList(t, db, "Work Tasks", entity.UserID(1))
				testutil.CreateTodoList(t, db, "Work Projects", entity.UserID(1))
				testutil.CreateTodoList(t, db, "Personal", entity.UserID(1))

				return db
			},
			opts: gatewayinput.ListTodoListsOptions{
				OwnerID:    testutil.UserIDPtr(entity.UserID(1)),
				NameSearch: testutil.StrPtr("Work"),
				Limit:      10,
			},
			wantLen:   2,
			wantTotal: 2,
		},

		"success: insensitive search": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				testutil.CreateTodoList(t, db, "WORK TASKS", entity.UserID(1))
				testutil.CreateTodoList(t, db, "work tasks", entity.UserID(1))

				return db
			},
			opts: gatewayinput.ListTodoListsOptions{
				OwnerID:    testutil.UserIDPtr(entity.UserID(1)),
				NameSearch: testutil.StrPtr("work"),
				Limit:      10,
			},
			wantLen:   2,
			wantTotal: 2,
		},

		"success: pagination": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				for i := 0; i < 5; i++ {
					testutil.CreateTodoList(
						t,
						db,
						fmt.Sprintf("List %d", i),
						entity.UserID(1),
					)
				}

				return db
			},
			opts: gatewayinput.ListTodoListsOptions{
				OwnerID: testutil.UserIDPtr(entity.UserID(1)),
				Limit:   2,
				Offset:  2,
			},
			wantLen:   2,
			wantTotal: 5,
		},

		"success: empty": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				return testutil.NewTestDB(t, cfg)
			},
			opts: gatewayinput.ListTodoListsOptions{
				OwnerID: testutil.UserIDPtr(entity.UserID(1)),
				Limit:   10,
			},
			wantLen:   0,
			wantTotal: 0,
		},

		"success: list by assignee": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				list1 := testutil.CreateTodoList(t, db, "Assigned List 1", entity.UserID(1))
				list2 := testutil.CreateTodoList(t, db, "Assigned List 2", entity.UserID(2))
				list3 := testutil.CreateTodoList(t, db, "Not Assigned", entity.UserID(3))

				testutil.CreateTodoWithAssignee(t, db, list1.ID, "Task 1", entity.UserID(10))
				testutil.CreateTodoWithAssignee(t, db, list2.ID, "Task 2", entity.UserID(10))
				testutil.CreateTodoWithAssignee(t, db, list3.ID, "Task 3", entity.UserID(99))

				return db
			},
			opts: gatewayinput.ListTodoListsOptions{
				AssigneeID: testutil.UserIDPtr(entity.UserID(10)),
				Limit:      10,
			},
			wantLen:   2,
			wantTotal: 2,
		},

		"success: list all by owner or assignee without duplicates": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				ownedAndAssigned := testutil.CreateTodoList(t, db, "Owned And Assigned", entity.UserID(1))
				ownedOnly := testutil.CreateTodoList(t, db, "Owned Only", entity.UserID(1))
				assignedOnly := testutil.CreateTodoList(t, db, "Assigned Only", entity.UserID(2))
				other := testutil.CreateTodoList(t, db, "Other", entity.UserID(3))

				testutil.CreateTodoWithAssignee(t, db, ownedAndAssigned.ID, "Task A", entity.UserID(1))
				testutil.CreateTodoWithAssignee(t, db, ownedAndAssigned.ID, "Task B", entity.UserID(1))
				testutil.CreateTodo(t, db, ownedOnly.ID, "Task C")
				testutil.CreateTodoWithAssignee(t, db, assignedOnly.ID, "Task D", entity.UserID(1))
				testutil.CreateTodoWithAssignee(t, db, other.ID, "Task E", entity.UserID(99))

				return db
			},
			opts: gatewayinput.ListTodoListsOptions{
				OwnerID:    testutil.UserIDPtr(entity.UserID(1)),
				AssigneeID: testutil.UserIDPtr(entity.UserID(1)),
				Limit:      10,
			},
			wantLen:   3,
			wantTotal: 3,
		},

		"success: list all with name search": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				ownedMatch := testutil.CreateTodoList(t, db, "Work Owned", entity.UserID(1))
				assignedMatch := testutil.CreateTodoList(t, db, "Work Assigned", entity.UserID(2))
				notMatch := testutil.CreateTodoList(t, db, "Personal", entity.UserID(1))

				testutil.CreateTodo(t, db, ownedMatch.ID, "Task A")
				testutil.CreateTodoWithAssignee(t, db, assignedMatch.ID, "Task B", entity.UserID(1))
				testutil.CreateTodo(t, db, notMatch.ID, "Task C")

				return db
			},
			opts: gatewayinput.ListTodoListsOptions{
				OwnerID:    testutil.UserIDPtr(entity.UserID(1)),
				AssigneeID: testutil.UserIDPtr(entity.UserID(1)),
				NameSearch: testutil.StrPtr("work"),
				Limit:      10,
			},
			wantLen:   2,
			wantTotal: 2,
		},

		"success: soft-deleted assigned todos are excluded": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				list := testutil.CreateTodoList(t, db, "List With Deleted Assigned Todo", entity.UserID(1))
				todo := testutil.CreateTodoWithAssignee(t, db, list.ID, "Deleted Assigned Todo", entity.UserID(10))

				cmdRepo := persistence.NewTodoCommandsGateway(db)

				err := cmdRepo.Delete(context.Background(), todo.ID)
				require.NoError(t, err)

				return db
			},
			opts: gatewayinput.ListTodoListsOptions{
				AssigneeID: testutil.UserIDPtr(entity.UserID(10)),
				Limit:      10,
			},
			wantLen:   0,
			wantTotal: 0,
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := tt.setup(t)

			repo := persistence.NewTodoListQueriesGateway(db)

			got, total, err := repo.List(
				context.Background(),
				&tt.opts,
			)

			require.NoError(t, err)

			assert.Len(t, got, tt.wantLen)
			assert.Equal(t, tt.wantTotal, total)
		})
	}
}
