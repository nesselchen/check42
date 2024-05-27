package stores

import (
	"check42/model/todos"
	"errors"
)

type TodoStore interface {
	GetAllTodos() ([]todos.Todo, error)
	UpdateTodo(int, todos.Todo) error
	GetTodo(int) (todos.Todo, error)
	CreateTodo(todos.Todo) error
	DeleteTodo(int) error
}

var ErrNotFound = errors.New("item not found")
