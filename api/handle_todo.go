package api

import (
	"check42/api/router"
	"check42/model"
	"check42/store/stores"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// GET /api/todo
func (s server) handleGetTodos(r *http.Request) ([]model.Todo, router.HttpStatus) {
	claims, ok := router.GetClaims(r)
	if !ok {
		return nil, router.HttpStatus{Code: http.StatusUnauthorized, Err: errors.New("insufficient claims")}
	}
	ts, err := s.todos.GetAllTodos(claims.ID)
	if err != nil {
		return nil, router.HttpStatus{Code: http.StatusInternalServerError, Err: err}
	}
	return ts, router.HttpStatus{Code: 200, Err: nil}
}

// POST /api/todo
func (s server) handlePostTodo(r *http.Request) (int64, router.HttpStatus) {
	claims, ok := router.GetClaims(r)
	if !ok {
		return 0, internalError
	}

	var todo model.CreateTodo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		return 0, router.HttpStatus{Code: http.StatusBadRequest, Err: err}
	}
	if err := todo.ValidateNew(); err.Err() {
		return 0, router.HttpStatus{Code: http.StatusBadRequest, Err: err}
	}

	todo.Owner = claims.ID
	id, err := s.todos.CreateTodo(todo)
	if err != nil {
		return 0, router.HttpStatus{Code: 500, Err: err}
	}
	return id, statusCreated
}

// GET /api/todo/{id}
func (s server) handleGetTodo(r *http.Request) (model.Todo, router.HttpStatus) {
	claims, ok := router.GetClaims(r)
	if !ok {
		return model.Todo{}, internalError
	}

	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		return model.Todo{}, badRequestCause(err)
	}
	td, err := s.todos.GetTodo(id, claims.ID)
	if err == stores.ErrNotFound {
		return model.Todo{}, notFound(id)
	}
	if err != nil {
		return model.Todo{}, internalErrorCause(err)
	}
	return td, statusOK
}

// DELETE /api/todo/{id}
func (s server) handleDeleteTodo(r *http.Request) router.HttpStatus {
	claims, ok := router.GetClaims(r)
	if !ok {
		return internalError
	}

	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		return badRequestCause(err)
	}
	err = s.todos.DeleteTodo(id, claims.ID)
	if err != nil {
		return internalErrorCause(err)
	}
	return router.HttpStatus{Code: http.StatusOK, Err: nil}
}

// PUT /api/todo/{id}
func (s server) handlePutTodo(r *http.Request) router.HttpStatus {
	claims, ok := router.GetClaims(r)
	if !ok {
		return internalError
	}

	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		return badRequestCause(err)
	}
	var t model.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return badRequestCause(err)
	}

	if err := s.todos.UpdateTodo(id, claims.ID, t); err != nil {
		return internalErrorCause(err)
	}
	return router.HttpStatus{Code: http.StatusOK, Err: nil}
}

// Updates the fields provided in the URL parameters.
// Options are done={bool} and text={string}.
// All other fields are preserved.
//
// PATCH /api/todo/{id}
func (s server) handlePatchTodo(r *http.Request) router.HttpStatus {
	claims, ok := router.GetClaims(r)
	if !ok {
		return internalError
	}

	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		return badRequestCause(err)
	}

	todo, err := s.todos.GetTodo(id, claims.ID)
	if err != nil {
		return internalError
	}
	if val := r.URL.Query().Get("done"); val != "" {
		done, err := strconv.ParseBool(val)
		if err != nil {
			return badRequestCause(err)
		}
		todo.Done = done
	}
	if val := r.URL.Query().Get("text"); val != "" {
		todo.Text = val
	}
	s.todos.UpdateTodo(todo.ID, todo.Owner, todo)
	return router.HttpStatus{Code: http.StatusOK, Err: nil}
}
