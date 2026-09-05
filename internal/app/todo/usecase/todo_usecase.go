package usecase

import (
	"errors"
	"strings"

	"github.com/jevvonn/ws-cicd-bcc2026/internal/app/todo/entity"
	"github.com/jevvonn/ws-cicd-bcc2026/internal/app/todo/repository"
	"gorm.io/gorm"
)

var (
	ErrTodoNotFound  = errors.New("todo not found")
	ErrTitleRequired = errors.New("title is required")
)

type TodoUsecase interface {
	Create(req entity.CreateTodoRequest) (*entity.Todo, error)
	GetAll() ([]entity.Todo, error)
	GetByID(id uint) (*entity.Todo, error)
	Update(id uint, req entity.UpdateTodoRequest) (*entity.Todo, error)
	Delete(id uint) error
}

type todoUsecase struct {
	repo repository.TodoRepository
}

func NewTodoUsecase(repo repository.TodoRepository) TodoUsecase {
	return &todoUsecase{repo: repo}
}

func (u *todoUsecase) Create(req entity.CreateTodoRequest) (*entity.Todo, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, ErrTitleRequired
	}

	todo := &entity.Todo{
		Title:       title,
		Description: req.Description,
		Completed:   false,
	}

	if err := u.repo.Create(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (u *todoUsecase) GetAll() ([]entity.Todo, error) {
	return u.repo.FindAll()
}

func (u *todoUsecase) GetByID(id uint) (*entity.Todo, error) {
	todo, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTodoNotFound
		}
		return nil, err
	}
	return todo, nil
}

func (u *todoUsecase) Update(id uint, req entity.UpdateTodoRequest) (*entity.Todo, error) {
	todo, err := u.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title == "" {
			return nil, ErrTitleRequired
		}
		todo.Title = title
	}

	if req.Description != nil {
		todo.Description = *req.Description
	}

	if req.Completed != nil {
		todo.Completed = *req.Completed
	}

	if err := u.repo.Update(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (u *todoUsecase) Delete(id uint) error {
	todo, err := u.GetByID(id)
	if err != nil {
		return err
	}
	return u.repo.Delete(todo)
}
