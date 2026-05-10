package graph

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	authinput "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/auth/input"
	authoutput "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/auth/output"
	todoinput "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/todo/input"
	todooutput "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/todo/output"
)

type TodoCreatorUsecase interface {
	CreateTodoList(ctx context.Context, in *todoinput.CreateTodoListInput) (*entity.TodoList, error)
	CreateTodo(ctx context.Context, in *todoinput.CreateTodoInput) (*entity.Todo, error)
}

type TodoGetterUsecase interface {
	GetTodoList(ctx context.Context, name string) (*entity.TodoList, error)
	GetTodo(ctx context.Context, name string) (*entity.Todo, error)
}

type TodoListerUsecase interface {
	ListTodoLists(ctx context.Context, parent string, opts *todoinput.ListTodoListsOptions) (*todooutput.TodoListPage, error)
	ListTodos(ctx context.Context, parent string, opts *todoinput.ListTodosOptions) (*todooutput.TodoPage, error)
}

type TodoUpdaterUsecase interface {
	UpdateTodoList(ctx context.Context, in *todoinput.UpdateTodoListInput) (*entity.TodoList, error)
	UpdateTodo(ctx context.Context, in *todoinput.UpdateTodoInput) (*entity.Todo, error)
}

type TodoDeleterUsecase interface {
	DeleteTodoList(ctx context.Context, name string) error
	DeleteTodo(ctx context.Context, name string) error
}

type AuthRegistererUsecase interface {
	Register(ctx context.Context, in *authinput.RegisterInput) (*entity.User, error)
}

type AuthLoginerUsecase interface {
	Login(ctx context.Context, in *authinput.LoginInput) (*authoutput.LoginOutput, error)
}

type UserGetterUsecase interface {
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByIDs(ctx context.Context, ids []string) ([]*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

type Resolver struct {
	todoCreator    TodoCreatorUsecase
	todoGetter     TodoGetterUsecase
	todoLister     TodoListerUsecase
	todoUpdater    TodoUpdaterUsecase
	todoDeleter    TodoDeleterUsecase
	authRegisterer AuthRegistererUsecase
	authLoginer    AuthLoginerUsecase
	userGetter     UserGetterUsecase
}

func NewResolver(
	todoCreator TodoCreatorUsecase,
	todoGetter TodoGetterUsecase,
	todoLister TodoListerUsecase,
	todoUpdater TodoUpdaterUsecase,
	todoDeleter TodoDeleterUsecase,
	authRegisterer AuthRegistererUsecase,
	authLoginer AuthLoginerUsecase,
	userGetter UserGetterUsecase,
) *Resolver {
	return &Resolver{
		todoCreator:    todoCreator,
		todoGetter:     todoGetter,
		todoLister:     todoLister,
		todoUpdater:    todoUpdater,
		todoDeleter:    todoDeleter,
		authRegisterer: authRegisterer,
		authLoginer:    authLoginer,
		userGetter:     userGetter,
	}
}

func (r *Resolver) GetUserGetter() UserGetterUsecase {
	return r.userGetter
}
