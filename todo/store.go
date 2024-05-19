package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"time"
)

type TodoStore interface {
	GetAllTodos() ([]Todo, error)
	GetTodo(int) (Todo, error)
	CreateTodo(Todo) error
	DeleteTodo(int) error
}

type jsonTodoStore struct {
	filename string
}

var ErrNotFound = errors.New("item not found")

func NewJsonTodoStore(filename string) (jsonTodoStore, error) {
	file, err := os.OpenFile(filename, os.O_CREATE, 0666)
	if err != nil {
		panic("Could neither create nor open JSON database file")
	}
	defer file.Close()
	return jsonTodoStore{filename}, nil
}

func appendTodo(todos []Todo, todo Todo) []Todo {
	if len(todos) == 0 {
		todo.Id = 0
		return append(todos, todo)
	}
	todo.Id = 1 + slices.MaxFunc(todos, func(t1, t2 Todo) int {
		return t1.Id - t2.Id
	}).Id
	return append(todos, todo)
}

func (store jsonTodoStore) CreateTodo(todo Todo) error {
	todos, err := store.read()
	if err != nil {
		return err
	}

	todo.Created = time.Now()
	todos = appendTodo(todos, todo)

	return store.write(todos)
}

func (store jsonTodoStore) GetAllTodos() ([]Todo, error) {
	return store.read()
}

func (store jsonTodoStore) DeleteTodo(id int) error {
	todos, err := store.read()
	if err != nil {
		return err
	}
	idx := slices.IndexFunc(todos, func(t Todo) bool {
		return t.Id == id
	})
	if idx == -1 {
		return ErrNotFound
	}
	todos[idx], todos[len(todos)-1] = todos[len(todos)-1], todos[idx]
	todos = todos[:len(todos)-1]
	store.write(todos)
	return nil
}

func (store jsonTodoStore) GetTodo(id int) (Todo, error) {
	todos, err := store.read()
	if err != nil {
		return Todo{}, err
	}
	idx := slices.IndexFunc(todos, func(t Todo) bool { return t.Id == id })
	if idx < 0 {
		return Todo{}, fmt.Errorf("no element with ID %d", id)
	}
	return todos[idx], nil
}

func (store jsonTodoStore) read() ([]Todo, error) {
	f, err := os.Open(store.filename)
	if err != nil {
		return nil, err
	}
	var todos []Todo
	err = json.NewDecoder(f).Decode(&todos)
	if err != nil {
		return nil, err
	}
	f.Close()
	return todos, nil
}

func (store jsonTodoStore) write(todos []Todo) error {
	f, err := os.OpenFile(store.filename, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(&todos)
}
