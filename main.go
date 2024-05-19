package main

import (
	"check42/router"
	"check42/todo"
	"encoding/json"
	"net/http"
	"strconv"
)

func main() {
	store := todo.NewJsonTodoStore("db.json")
	th := todoHandler{store}

	todoRouter := router.NewMethodMux("/todo")
	todoRouter.OnGet(router.Process(th.handleGetTodos))
	todoRouter.OnPost(router.ProcessWithoutResponseBody(th.handlePostTodo))

	singleTodoRouter := router.NewMethodMux("/todo/{id}")
	singleTodoRouter.OnGet(router.Process(th.handleGetTodo))
	singleTodoRouter.OnDelete(router.ProcessWithoutResponseBody(th.handleDeleteTodo))

	router.ListenAndServe("127.0.0.1:9999", todoRouter, singleTodoRouter)
}

type todoHandler struct {
	store todo.TodoStore
}

// GET /todo
func (th todoHandler) handleGetTodos(r *http.Request) ([]todo.Todo, router.HttpStatus) {
	ts, err := th.store.GetAllTodos()
	if err != nil {
		return nil, router.HttpStatus{Code: http.StatusInternalServerError, Err: err}
	}
	return ts, router.HttpStatus{Code: 200, Err: nil}
}

// POST /todo
func (th todoHandler) handlePostTodo(r *http.Request) router.HttpStatus {
	var todo todo.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		return router.HttpStatus{Code: http.StatusBadRequest, Err: err}
	}
	err = th.store.CreateTodo(todo)
	if err != nil {
		return router.HttpStatus{Code: 500, Err: err}
	}
	return router.HttpStatus{Code: 201, Err: nil}
}

// GET /todo/{id}
func (th todoHandler) handleGetTodo(r *http.Request) (todo.Todo, router.HttpStatus) {
	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 32)
	if err != nil {
		return todo.Todo{}, router.HttpStatus{Code: http.StatusBadRequest, Err: err}
	}
	td, err := th.store.GetTodo(int(id))
	if err != nil {
		return todo.Todo{}, router.HttpStatus{Code: http.StatusNotFound, Err: nil}
	}
	return td, router.HttpStatus{Code: http.StatusOK, Err: nil}
}

// DELETE /todo/{id}
func (th todoHandler) handleDeleteTodo(r *http.Request) router.HttpStatus {
	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 32)
	if err != nil {
		return router.HttpStatus{Code: http.StatusBadRequest, Err: err}
	}
	err = th.store.DeleteTodo(int(id))
	if err != nil {
		return router.HttpStatus{Code: http.StatusInternalServerError, Err: err}
	}
	return router.HttpStatus{Code: http.StatusOK, Err: nil}
}
