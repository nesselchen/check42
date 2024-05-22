package store

import (
	"check42/model/todos"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"time"
)

type jsonTodoStore struct {
	filename string
}

func NewJsonTodoStore(filename string) jsonTodoStore {
	return jsonTodoStore{filename}
}

func appendTodo(tds []todos.Todo, todo todos.Todo) []todos.Todo {
	if len(tds) == 0 {
		todo.Id = 0
		return append(tds, todo)
	}
	todo.Id = 1 + slices.MaxFunc(tds, func(t1, t2 todos.Todo) int {
		return t1.Id - t2.Id
	}).Id
	return append(tds, todo)
}

func (store jsonTodoStore) CreateTodo(todo todos.Todo) error {
	todos, err := store.read()
	if err != nil {
		return err
	}

	todo.Created = time.Now()
	todos = appendTodo(todos, todo)

	return store.write(todos)
}

func (store jsonTodoStore) GetAllTodos() ([]todos.Todo, error) {
	return store.read()
}

func (store jsonTodoStore) DeleteTodo(id int) error {
	tds, err := store.read()
	if err != nil {
		return err
	}
	idx := slices.IndexFunc(tds, func(t todos.Todo) bool {
		return t.Id == id
	})
	if idx == -1 {
		return ErrNotFound
	}
	// swap with last and decrement length
	tds[idx], tds[len(tds)-1] = tds[len(tds)-1], tds[idx]
	tds = tds[:len(tds)-1]
	store.write(tds)
	return nil
}

func (store jsonTodoStore) GetTodo(id int) (todos.Todo, error) {
	tds, err := store.read()
	if err != nil {
		return todos.Todo{}, err
	}
	idx := slices.IndexFunc(tds, func(t todos.Todo) bool { return t.Id == id })
	if idx < 0 {
		return todos.Todo{}, fmt.Errorf("no element with ID %d", id)
	}
	return tds[idx], nil
}

func (store jsonTodoStore) read() ([]todos.Todo, error) {
	f, err := os.OpenFile(store.filename, os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tds []todos.Todo
	err = json.NewDecoder(f).Decode(&tds)

	if err == io.EOF {
		return make([]todos.Todo, 0), nil
	} else if err != nil {
		return nil, err
	}

	return tds, nil
}

func (store jsonTodoStore) write(todos []todos.Todo) error {
	f, err := os.OpenFile(store.filename, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(&todos)
}
