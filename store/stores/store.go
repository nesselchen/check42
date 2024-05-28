package stores

import (
	"check42/model"
	"errors"
)

type UserStore interface {
	GetUserByID(id int) (model.User, error)
	GetUserByName(name string) (model.User, error)
	CreateUser(model.NewUser) error
}

type TodoStore interface {
	GetAllTodos(userID int64) ([]model.Todo, error)
	UpdateTodo(todoID, userID int64, update model.Todo) error
	GetTodo(todoID, userID int64) (model.Todo, error)
	CreateTodo(model.Todo) (int64, error)
	DeleteTodo(todoID, userID int64) error
}

var (
	ErrNotFound      = errors.New("item not found")
	ErrUsernameTaken = errors.New("username is taken")
	ErrEmailTaken    = errors.New("email is taken")
)