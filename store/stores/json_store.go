package stores

import (
	"check42/model"
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

func appendTodo(tds []model.Todo, todo model.Todo) []model.Todo {
	if len(tds) == 0 {
		todo.ID = 0
		return append(tds, todo)
	}
	todo.ID = 1 + slices.MaxFunc(tds, func(t1, t2 model.Todo) int {
		return t1.ID - t2.ID
	}).ID
	return append(tds, todo)
}

func (store jsonTodoStore) CreateTodo(todo model.Todo) error {
	todos, err := store.read()
	if err != nil {
		return err
	}

	todo.Created = time.Now()
	todos = appendTodo(todos, todo)

	return store.write(todos)
}

func (store jsonTodoStore) GetAllTodos() ([]model.Todo, error) {
	return store.read()
}

func (store jsonTodoStore) DeleteTodo(id int) error {
	tds, err := store.read()
	if err != nil {
		return err
	}
	idx := slices.IndexFunc(tds, func(t model.Todo) bool {
		return t.ID == id
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

func (store jsonTodoStore) GetTodo(id int) (model.Todo, error) {
	tds, err := store.read()
	if err != nil {
		return model.Todo{}, err
	}
	idx := slices.IndexFunc(tds, func(t model.Todo) bool { return t.ID == id })
	if idx < 0 {
		return model.Todo{}, fmt.Errorf("no element with ID %d", id)
	}
	return tds[idx], nil
}

func (store jsonTodoStore) UpdateTodo(id int, t model.Todo) error {
	panic("Calling unimplemented method jsonTodoStore.UpdateTodo")
}

func (store jsonTodoStore) read() ([]model.Todo, error) {
	f, err := os.OpenFile(store.filename, os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tds []model.Todo
	err = json.NewDecoder(f).Decode(&tds)

	if err == io.EOF {
		return make([]model.Todo, 0), nil
	} else if err != nil {
		return nil, err
	}

	return tds, nil
}

func (store jsonTodoStore) write(todos []model.Todo) error {
	f, err := os.OpenFile(store.filename, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(&todos)
}
