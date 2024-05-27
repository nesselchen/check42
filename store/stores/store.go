package stores

import (
	"check42/model"
	"errors"
)

type TodoStore interface {
	GetAllTodos() ([]model.Todo, error)
	UpdateTodo(int, model.Todo) error
	GetTodo(int) (model.Todo, error)
	CreateTodo(model.Todo) error
	DeleteTodo(int) error
}

var ErrNotFound = errors.New("item not found")
