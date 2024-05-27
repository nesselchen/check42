package api

import (
	"check42/api/router"
	"check42/model/todos"
	"check42/store/stores"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

type server struct {
	addr  string
	store stores.TodoStore
}

func RunServer(addr string, store stores.TodoStore) {
	s := &server{addr, store}

	baseRoute := router.New("/")
	apiRoute := baseRoute.Subroute("api")
	todoRoute := apiRoute.Subroute("/todo")
	todoIdRoute := todoRoute.Subroute("/{id}")

	baseRoute.OnGet(s.templateHtml)

	todoRoute.OnPost(router.ProcessWithoutResponseBody(s.handlePostTodo))
	todoRoute.OnGet(router.Process(s.handleGetTodos))

	todoIdRoute.OnGet(router.Process(s.handleGetTodo))
	todoIdRoute.OnDelete(router.ProcessWithoutResponseBody(s.handleDeleteTodo))
	todoIdRoute.OnPatch(router.ProcessWithoutResponseBody(s.handlePatchTodo))

	log.Fatal(router.ListenAndServe(s.addr, baseRoute))
}

// GET /
func (s server) templateHtml(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index.html.tmpl").ParseFiles("templates/index.html.tmpl")
	if err != nil {
		w.WriteHeader(500)
		return
	}
	todos, err := s.store.GetAllTodos()
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
func (s server) handleGetTodos(r *http.Request) ([]todos.Todo, router.HttpStatus) {
	ts, err := s.store.GetAllTodos()
	if err != nil {
		return nil, router.HttpStatus{Code: http.StatusInternalServerError, Err: err}
	}
	return ts, router.HttpStatus{Code: 200, Err: nil}
}

// POST /api/todo
func (s server) handlePostTodo(r *http.Request) router.HttpStatus {
	var todo todos.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		return router.HttpStatus{Code: http.StatusBadRequest, Err: err}
	}
	if err := todo.ValidateNew(); err != nil {
		return router.HttpStatus{Code: http.StatusBadRequest, Err: err}
	}
	err = s.store.CreateTodo(todo)
	if err != nil {
		return router.HttpStatus{Code: 500, Err: err}
	}
	return router.HttpStatus{Code: 201, Err: nil}
}

// GET /api/todo/{id}
func (s server) handleGetTodo(r *http.Request) (todos.Todo, router.HttpStatus) {
	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		return todos.Todo{}, router.HttpStatus{Code: http.StatusBadRequest, Err: err}
	}
	td, err := s.store.GetTodo(int(id))
	if err == stores.ErrNotFound {
		return todos.Todo{}, router.HttpStatus{Code: http.StatusNotFound, Err: fmt.Errorf("no todo with ID %d", id)}
	}
	if err != nil {
		return todos.Todo{}, router.HttpStatus{Code: http.StatusInternalServerError}
	}
	return td, router.HttpStatus{Code: http.StatusOK, Err: nil}
}

// DELETE /api/todo/{id}
func (s server) handleDeleteTodo(r *http.Request) router.HttpStatus {
	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		return router.HttpStatus{Code: http.StatusBadRequest, Err: err}
	}
	err = s.store.DeleteTodo(int(id))
	if err != nil {
		return router.HttpStatus{Code: http.StatusInternalServerError, Err: err}
	}
	return router.HttpStatus{Code: http.StatusOK, Err: nil}
}

// PATCH /api/todo/{id}
func (s server) handlePatchTodo(r *http.Request) router.HttpStatus {
	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		return router.HttpStatus{Code: http.StatusBadRequest, Err: err}
	}
	var t todos.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return router.HttpStatus{Code: http.StatusBadRequest, Err: err}
	}
	if err := s.store.UpdateTodo(int(id), t); err != nil {
		return router.HttpStatus{Code: http.StatusInternalServerError, Err: err}
	}
	return router.HttpStatus{Code: http.StatusOK, Err: nil}
}
