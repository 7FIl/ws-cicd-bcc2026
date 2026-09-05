package repository

import (
	"github.com/jevvonn/ws-cicd-bcc2026/internal/app/todo/entity"
	"gorm.io/gorm"
)

type TodoRepository interface {
	Create(todo *entity.Todo) error
	FindAll() ([]entity.Todo, error)
	FindByID(id uint) (*entity.Todo, error)
	Update(todo *entity.Todo) error
	Delete(todo *entity.Todo) error
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(todo *entity.Todo) error {
	return r.db.Create(todo).Error
}

func (r *todoRepository) FindAll() ([]entity.Todo, error) {
	todos := []entity.Todo{}
	err := r.db.Order("id asc").Find(&todos).Error
	return todos, err
}

func (r *todoRepository) FindByID(id uint) (*entity.Todo, error) {
	var todo entity.Todo
	if err := r.db.First(&todo, id).Error; err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepository) Update(todo *entity.Todo) error {
	return r.db.Save(todo).Error
}

func (r *todoRepository) Delete(todo *entity.Todo) error {
	return r.db.Delete(todo).Error
}
