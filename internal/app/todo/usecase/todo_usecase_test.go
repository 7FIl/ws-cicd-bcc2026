package usecase_test

import (
	"testing"

	"github.com/jevvonn/ws-cicd-bcc2026/internal/app/todo/entity"
	"github.com/jevvonn/ws-cicd-bcc2026/internal/app/todo/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type mockTodoRepository struct {
	mock.Mock
}

func (m *mockTodoRepository) Create(todo *entity.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *mockTodoRepository) FindAll() ([]entity.Todo, error) {
	args := m.Called()
	return args.Get(0).([]entity.Todo), args.Error(1)
}

func (m *mockTodoRepository) FindByID(id uint) (*entity.Todo, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Todo), args.Error(1)
}

func (m *mockTodoRepository) Update(todo *entity.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *mockTodoRepository) Delete(todo *entity.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func TestCreateTodo(t *testing.T) {
	repo := new(mockTodoRepository)
	repo.On("Create", mock.AnythingOfType("*entity.Todo")).Return(nil)

	todoUsecase := usecase.NewTodoUsecase(repo)
	todo, err := todoUsecase.Create(entity.CreateTodoRequest{
		Title:       "Learn Go",
		Description: "Build a Fiber API",
	})

	assert.NoError(t, err)
	assert.Equal(t, "Learn Go", todo.Title)
	assert.False(t, todo.Completed)
	repo.AssertExpectations(t)
}

func TestCreateTodoEmptyTitle(t *testing.T) {
	repo := new(mockTodoRepository)
	todoUsecase := usecase.NewTodoUsecase(repo)

	todo, err := todoUsecase.Create(entity.CreateTodoRequest{Title: "   "})

	assert.Nil(t, todo)
	assert.ErrorIs(t, err, usecase.ErrTitleRequired)
	repo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestGetAllTodo(t *testing.T) {
	repo := new(mockTodoRepository)
	repo.On("FindAll").Return([]entity.Todo{{ID: 1, Title: "Learn Go"}}, nil)

	todoUsecase := usecase.NewTodoUsecase(repo)
	todos, err := todoUsecase.GetAll()

	assert.NoError(t, err)
	assert.Len(t, todos, 1)
	repo.AssertExpectations(t)
}

func TestGetTodoByIDNotFound(t *testing.T) {
	repo := new(mockTodoRepository)
	repo.On("FindByID", uint(99)).Return(nil, gorm.ErrRecordNotFound)

	todoUsecase := usecase.NewTodoUsecase(repo)
	todo, err := todoUsecase.GetByID(99)

	assert.Nil(t, todo)
	assert.ErrorIs(t, err, usecase.ErrTodoNotFound)
	repo.AssertExpectations(t)
}

func TestUpdateTodo(t *testing.T) {
	completed := true
	existing := &entity.Todo{ID: 1, Title: "Learn Go", Completed: false}

	repo := new(mockTodoRepository)
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	todoUsecase := usecase.NewTodoUsecase(repo)
	todo, err := todoUsecase.Update(1, entity.UpdateTodoRequest{Completed: &completed})

	assert.NoError(t, err)
	assert.True(t, todo.Completed)
	repo.AssertExpectations(t)
}

func TestDeleteTodo(t *testing.T) {
	existing := &entity.Todo{ID: 1, Title: "Learn Go"}

	repo := new(mockTodoRepository)
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Delete", existing).Return(nil)

	todoUsecase := usecase.NewTodoUsecase(repo)

	assert.NoError(t, todoUsecase.Delete(1))
	repo.AssertExpectations(t)
}
