package api

import (
	"check42/api/router"
	"check42/model"
	"errors"
	"net/http"
	"strconv"
)

// GET /api/todo/category
func (s server) handleGetCategories(r *http.Request) ([]model.TodoCategory, router.HttpStatus) {
	claims, ok := router.GetClaims(r)
	if !ok {
		return nil, internalError
	}

	cats, err := s.todos.GetAllCategories(claims.ID)
	if err != nil {
		return nil, internalErrorCause(err)
	}

	return cats, statusOK
}

// POST /api/todo/category?name={name}
func (s server) handlePostCategory(r *http.Request) (int64, router.HttpStatus) {
	claims, ok := router.GetClaims(r)
	if !ok {
		return 0, internalError
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		return 0, badRequestCause(errors.New("missing field 'name'"))
	}

	id, err := s.todos.CreateCategory(name, claims.ID)
	if err != nil {
		return 0, internalErrorCause(err)
	}

	return id, statusOK
}

// PATCH /api/todo/category/{id}?name={name}
func (s server) handlePatchCategory(r *http.Request) router.HttpStatus {
	claims, ok := router.GetClaims(r)
	if !ok {
		return internalError
	}

	pathValue := r.PathValue("id")
	categoryID, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		return badRequestCause(errors.New("incorrect 'id'"))
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		return badRequestCause(errors.New("missing field 'name'"))
	}

	err = s.todos.UpdateCategory(name, categoryID, claims.ID)
	if err != nil {
		return internalErrorCause(err)
	}

	return statusOK
}

// DELETE /api/todo/category/{id}
func (s server) handleDeleteCategory(r *http.Request) router.HttpStatus {
	claims, ok := router.GetClaims(r)
	if !ok {
		return internalError
	}

	pathValue := r.PathValue("id")
	categoryID, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		return badRequestCause(errors.New("incorrect 'id'"))
	}
	if categoryID == 0 {
		return badRequestCause(errors.New("cannot delete this category"))
	}

	err = s.todos.DeleteCategory(categoryID, claims.ID)
	if err != nil {
		return internalErrorCause(err)
	}

	return statusOK
}

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
