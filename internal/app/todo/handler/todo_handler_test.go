package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/ws-cicd-bcc2026/internal/app/todo/entity"
	"github.com/jevvonn/ws-cicd-bcc2026/internal/app/todo/handler"
	"github.com/jevvonn/ws-cicd-bcc2026/internal/app/todo/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTodoUsecase struct {
	mock.Mock
}

func (m *mockTodoUsecase) Create(req entity.CreateTodoRequest) (*entity.Todo, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Todo), args.Error(1)
}

func (m *mockTodoUsecase) GetAll() ([]entity.Todo, error) {
	args := m.Called()
	return args.Get(0).([]entity.Todo), args.Error(1)
}

func (m *mockTodoUsecase) GetByID(id uint) (*entity.Todo, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Todo), args.Error(1)
}

func (m *mockTodoUsecase) Update(id uint, req entity.UpdateTodoRequest) (*entity.Todo, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Todo), args.Error(1)
}

func (m *mockTodoUsecase) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func newTestApp(todoUsecase usecase.TodoUsecase) *fiber.App {
	app := fiber.New()
	handler.NewTodoHandler(app.Group("/api"), todoUsecase)
	return app
}

func TestCreateTodoHandler(t *testing.T) {
	todoUsecase := new(mockTodoUsecase)
	todoUsecase.On("Create", entity.CreateTodoRequest{Title: "Learn Go"}).
		Return(&entity.Todo{ID: 1, Title: "Learn Go"}, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/todos", strings.NewReader(`{"title":"Learn Go"}`))
	req.Header.Set("Content-Type", "application/json")

	res, err := newTestApp(todoUsecase).Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, res.StatusCode)
	todoUsecase.AssertExpectations(t)
}

func TestGetAllTodoHandler(t *testing.T) {
	todoUsecase := new(mockTodoUsecase)
	todoUsecase.On("GetAll").Return([]entity.Todo{{ID: 1, Title: "Learn Go"}}, nil)

	res, err := newTestApp(todoUsecase).Test(httptest.NewRequest(http.MethodGet, "/api/todos", nil))

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	var body map[string]any
	assert.NoError(t, json.NewDecoder(res.Body).Decode(&body))
	assert.Equal(t, true, body["success"])
	todoUsecase.AssertExpectations(t)
}

func TestGetTodoByIDNotFoundHandler(t *testing.T) {
	todoUsecase := new(mockTodoUsecase)
	todoUsecase.On("GetByID", uint(99)).Return(nil, usecase.ErrTodoNotFound)

	res, err := newTestApp(todoUsecase).Test(httptest.NewRequest(http.MethodGet, "/api/todos/99", nil))

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, res.StatusCode)
	todoUsecase.AssertExpectations(t)
}

func TestDeleteTodoHandler(t *testing.T) {
	todoUsecase := new(mockTodoUsecase)
	todoUsecase.On("Delete", uint(1)).Return(nil)

	res, err := newTestApp(todoUsecase).Test(httptest.NewRequest(http.MethodDelete, "/api/todos/1", nil))

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)
	todoUsecase.AssertExpectations(t)
}
