package todo_test

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/auth"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/grpc/todo"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

func authCtx(userID string) context.Context {
	return auth.SetUserContext(context.Background(), userID, []string{})
}

var (
	sampleTodoList = &entity.TodoList{
		ID:      entity.TodoListID(2),
		Name:    "Task 1",
		OwnerID: entity.UserID(1),
	}

	sampleTodo = &entity.Todo{
		ID:         entity.TodoID(3),
		TodoListID: entity.TodoListID(2),
		Title:      "Sub Task 1",
		Status:     entity.TodoStatusPending,
		Priority:   entity.PriorityMedium,
		CreatorID:  entity.UserID(1),
	}
)

// Todo stubs

type stubTodoCreator struct {
	gotInput *input.TodoCreator
	resp     *output.TodoCreator
	err      error
}

func (s *stubTodoCreator) Create(ctx context.Context, in *input.TodoCreator) (*output.TodoCreator, error) {
	s.gotInput = in
	return s.resp, s.err
}

type stubTodoGetter struct {
	gotInput *input.TodoGetter
	resp     *output.TodoGetter
	err      error
}

func (s *stubTodoGetter) Get(ctx context.Context, in *input.TodoGetter) (*output.TodoGetter, error) {
	s.gotInput = in
	return s.resp, s.err
}

type stubTodoLister struct {
	gotInput *input.TodoLister
	resp     *output.TodoLister
	err      error
}

func (s *stubTodoLister) List(ctx context.Context, in *input.TodoLister) (*output.TodoLister, error) {
	s.gotInput = in
	return s.resp, s.err
}

type stubTodoUpdater struct {
	gotInput *input.TodoUpdater
	resp     *output.TodoUpdater
	err      error
}

func (s *stubTodoUpdater) Update(ctx context.Context, in *input.TodoUpdater) (*output.TodoUpdater, error) {
	s.gotInput = in
	return s.resp, s.err
}

type stubTodoDeleter struct {
	gotInput *input.TodoDeleter
	err      error
}

func (s *stubTodoDeleter) Delete(ctx context.Context, in *input.TodoDeleter) (*output.TodoDeleter, error) {
	s.gotInput = in

	if s.err != nil {
		return nil, s.err
	}

	return &output.TodoDeleter{}, nil
}

// Todo List stubs

type stubTodoListCreator struct {
	gotInput *input.TodoListCreator
	resp     *output.TodoListCreator
	err      error
}

func (s *stubTodoListCreator) Create(ctx context.Context, in *input.TodoListCreator) (*output.TodoListCreator, error) {
	s.gotInput = in
	return s.resp, s.err
}

type stubTodoListGetter struct {
	gotInput *input.TodoListGetter
	resp     *output.TodoListGetter
	err      error
}

func (s *stubTodoListGetter) Get(ctx context.Context, in *input.TodoListGetter) (*output.TodoListGetter, error) {
	s.gotInput = in
	return s.resp, s.err
}

type stubTodoListLister struct {
	gotInput *input.TodoListLister
	resp     *output.TodoListLister
	err      error
}

func (s *stubTodoListLister) List(ctx context.Context, in *input.TodoListLister) (*output.TodoListLister, error) {
	s.gotInput = in
	return s.resp, s.err
}

type stubTodoListUpdater struct {
	gotInput *input.TodoListUpdater
	resp     *output.TodoListUpdater
	err      error
}

func (s *stubTodoListUpdater) Update(ctx context.Context, in *input.TodoListUpdater) (*output.TodoListUpdater, error) {
	s.gotInput = in
	return s.resp, s.err
}

type stubTodoListDeleter struct {
	gotInput *input.TodoListDeleter
	err      error
}

func (s *stubTodoListDeleter) Delete(ctx context.Context, in *input.TodoListDeleter) (*output.TodoListDeleter, error) {
	s.gotInput = in

	if s.err != nil {
		return nil, s.err
	}

	return &output.TodoListDeleter{}, nil
}

type handlerBuilder struct {
	todoCreator     todo.TodoCreatorUsecase
	todoGetter      todo.TodoGetterUsecase
	todoLister      todo.TodoListerUsecase
	todoUpdater     todo.TodoUpdaterUsecase
	todoDeleter     todo.TodoDeleterUsecase
	todoListCreator todo.TodoListCreatorUsecase
	todoListGetter  todo.TodoListGetterUsecase
	todoListLister  todo.TodoListListerUsecase
	todoListUpdater todo.TodoListUpdaterUsecase
	todoListDeleter todo.TodoListDeleterUsecase
}

func (b *handlerBuilder) build() *todo.TodoHandler {
	return todo.NewTodoHandler(
		b.todoCreator,
		b.todoGetter,
		b.todoLister,
		b.todoUpdater,
		b.todoDeleter,
		b.todoListCreator,
		b.todoListGetter,
		b.todoListLister,
		b.todoListUpdater,
		b.todoListDeleter,
	)
}
