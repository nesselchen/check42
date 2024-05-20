package main

import (
	"check42/router"
	"check42/todo"
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"
)

func main() {
	store := todo.NewJsonTodoStore("db.json")
	th := todoHandler{store}

	frontendRouter := router.NewRoute("/")
	frontendRouter.OnGet(th.templateHtml)

	todoRouter := router.NewRoute("/api/todo")
	todoRouter.OnGet(router.Process(th.handleGetTodos))
	todoRouter.OnPost(router.ProcessWithoutResponseBody(th.handlePostTodo))

	singleTodoRouter := router.NewRoute("/api/todo/{id}")
	singleTodoRouter.OnGet(router.Process(th.handleGetTodo))
	singleTodoRouter.OnDelete(router.ProcessWithoutResponseBody(th.handleDeleteTodo))

	router.ListenAndServe("0.0.0.0:2442", frontendRouter, todoRouter, singleTodoRouter)
}

type todoHandler struct {
	store todo.TodoStore
}

// GET /
func (th todoHandler) templateHtml(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index.html.tmpl").ParseFiles("templates/index.html.tmpl")
	if err != nil {
		w.WriteHeader(500)
		return
	}
	todos, err := th.store.GetAllTodos()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	if err := tmpl.Execute(w, todos); err != nil {
		w.WriteHeader(500)
		return
	}
}

// GET /api/todo
func (th todoHandler) handleGetTodos(r *http.Request) ([]todo.Todo, router.HttpStatus) {
	ts, err := th.store.GetAllTodos()
	if err != nil {
		return nil, router.HttpStatus{Code: http.StatusInternalServerError, Err: err}
	}
	return ts, router.HttpStatus{Code: 200, Err: nil}
}

// POST /api/todo
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

// GET /api/todo/{id}
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

// DELETE /api/todo/{id}
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
