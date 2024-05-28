package api

import (
	"check42/api/router"
	"check42/model"
	"check42/store/stores"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type server struct {
	addr  string
	todos stores.TodoStore
	users stores.UserStore
}

func RunServer(addr string, todos stores.TodoStore, users stores.UserStore) {
	s := &server{addr, todos, users}

	godotenv.Load()

	secret, found := os.LookupEnv("JWT_SECRET")
	if !found {
		log.Fatal("Fatal error: missing environment variable 'JWT_SECRET'")
	}
	authority := ApiAuthority{
		users,
		[]byte(secret),
	}

	login := router.New("/login")
	login.OnPost(s.handleLogin)
	login.Use(router.BasicAuth(authority))

	base := router.New("/")
	base.Use(router.JWTAuth(authority))

	api := base.Subroute("api")
	todo := api.Subroute("/todo")
	todoId := todo.Subroute("/{id}")

	user := api.Subroute("/user")
	user.OnGet(router.Process(s.handleGetUsers))

	// base.OnGet(s.templateHtml)

	todo.OnPost(router.ProcessWithoutResponseBody(s.handlePostTodo))
	todo.OnGet((router.Process(s.handleGetTodos)))

	todoId.OnGet(router.Process(s.handleGetTodo))
	todoId.OnDelete(router.ProcessWithoutResponseBody(s.handleDeleteTodo))
	todoId.OnPatch(router.ProcessWithoutResponseBody(s.handlePatchTodo))

	log.Fatal(router.ListenAndServe(s.addr, base, login))
}

// // GET /
// func (s server) templateHtml(w http.ResponseWriter, r *http.Request) {
// 	tmpl, err := template.New("index.html.tmpl").ParseFiles("templates/index.html.tmpl")
// 	if err != nil {
// 		w.WriteHeader(500)
// 		return
// 	}
// 	todos, err := s.todos.GetAllTodos()
// 	if err != nil {
// 		w.WriteHeader(500)
// 		return
// 	}
// 	if err := tmpl.Execute(w, todos); err != nil {
// 		w.WriteHeader(500)
// 		return
// 	}
// }

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
func (s server) handlePostTodo(r *http.Request) router.HttpStatus {
	var todo model.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		return router.HttpStatus{Code: http.StatusBadRequest, Err: err}
	}
	if err := todo.ValidateNew(); err.Err() {
		return router.HttpStatus{Code: http.StatusBadRequest, Err: err}
	}
	err = s.todos.CreateTodo(todo)
	if err != nil {
		return router.HttpStatus{Code: 500, Err: err}
	}
	return router.HttpStatus{Code: 201, Err: nil}
}

// GET /api/todo/{id}
func (s server) handleGetTodo(r *http.Request) (model.Todo, router.HttpStatus) {
	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		return model.Todo{}, router.HttpStatus{Code: http.StatusBadRequest, Err: err}
	}
	td, err := s.todos.GetTodo(int(id))
	if err == stores.ErrNotFound {
		return model.Todo{}, router.HttpStatus{Code: http.StatusNotFound, Err: fmt.Errorf("no todo with ID %d", id)}
	}
	if err != nil {
		return model.Todo{}, router.HttpStatus{Code: http.StatusInternalServerError}
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
	err = s.todos.DeleteTodo(int(id))
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
	var t model.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return router.HttpStatus{Code: http.StatusBadRequest, Err: err}
	}
	if err := s.todos.UpdateTodo(int(id), t); err != nil {
		return router.HttpStatus{Code: http.StatusInternalServerError, Err: err}
	}
	return router.HttpStatus{Code: http.StatusOK, Err: nil}
}

// GET /api/user
func (s server) handleGetUsers(r *http.Request) (model.User, router.HttpStatus) {
	u, err := s.users.GetUserByID(1)
	if err != nil {
		return model.User{}, router.HttpStatus{Code: http.StatusNotFound, Err: err}
	}
	return u, router.HttpStatus{Code: 200, Err: nil}
}

// POST /login
func (s server) handleLogin(w http.ResponseWriter, r *http.Request) {
	claims, ok := router.GetClaims(r)
	if !ok {
		fail(w, http.StatusInternalServerError, "internal error")
		return
	}
	week := time.Duration(7 * 24 * time.Hour)
	jwtClaims := jwt.MapClaims{
		"sub": claims.Name,
		"id":  claims.ID,
		"exp": jwt.NumericDate{Time: time.Now().Add(week)},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	secret, found := os.LookupEnv("JWT_SECRET")
	if !found {
		fail(w, http.StatusInternalServerError, "internal error")
		return
	}
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		fail(w, http.StatusInternalServerError, "internal error")
		return
	}

	io.WriteString(w, signed)
}

func fail(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	io.WriteString(w, msg)
}
