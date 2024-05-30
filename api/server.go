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
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type server struct {
	addr  string
	todos stores.TodoStore
	users stores.UserStore
}

func RunServer(addr string, todos stores.TodoStore, users stores.UserStore) {
	s := &server{addr, todos, users}

	secret, found := os.LookupEnv("JWT_SECRET")
	if !found {
		log.Fatal("Fatal error: missing environment variable 'JWT_SECRET'")
	}
	authority := ApiAuthority{
		store:     users,
		jwtSecret: []byte(secret),
		pwSalt:    os.Getenv("PW_SALT"),
	}

	// routes
	base := router.New("/")
	assets := base.Subroute("static/")

	auth := base.Subroute("auth")
	signin := auth.Subroute("/signin")
	login := auth.Subroute("/login")

	api := base.Subroute("api")
	todo := api.Subroute("/todo")
	todoId := todo.Subroute("/{id}")

	// middlewares
	base.Use(router.LogCall)
	login.Use(router.BasicAuth(authority))
	api.Use(router.JWTAuth(authority))

	// handlers
	signin.OnPost(s.handleSignin)

	login.OnPost(s.handleLogin)

	base.OnGet(handleBase)
	assets.OnGet(handleStatic)

	todo.OnPost(router.Process(s.handlePostTodo))
	todo.OnGet((router.Process(s.handleGetTodos)))

	todoId.OnGet(router.Process(s.handleGetTodo))
	todoId.OnDelete(router.ProcessWithoutResponseBody(s.handleDeleteTodo))
	todoId.OnPut(router.ProcessWithoutResponseBody(s.handlePutTodo))
	todoId.OnPatch(router.ProcessWithoutResponseBody(s.handlePatchTodo))

	log.Fatal(router.ListenAndServe(s.addr, base))
}

// GET /
func handleBase(w http.ResponseWriter, r *http.Request) {
	html, err := os.ReadFile("static/frontend/index.html")
	if err != nil {
		log.Fatal("Frontend files are missing")
	}
	io.WriteString(w, string(html))
}

// GET /static/{path}
func handleStatic(w http.ResponseWriter, r *http.Request) {
	path, ok := strings.CutPrefix(r.URL.Path, "/static/")
	file, err := os.ReadFile("./static/" + path)
	if !ok || err != nil {
		fail(w, 404, "file not found")
		return
	}
	io.WriteString(w, string(file))
}

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
		return 0, router.HttpStatus{Code: 500, Err: errors.New("internal error")}
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
	return id, router.HttpStatus{Code: 201, Err: nil}
}

// GET /api/todo/{id}
func (s server) handleGetTodo(r *http.Request) (model.Todo, router.HttpStatus) {
	claims, ok := router.GetClaims(r)
	if !ok {
		return model.Todo{}, router.HttpStatus{Code: 500, Err: errors.New("internal error")}
	}

	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		return model.Todo{}, router.HttpStatus{Code: http.StatusBadRequest, Err: err}
	}
	td, err := s.todos.GetTodo(id, claims.ID)
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
	claims, ok := router.GetClaims(r)
	if !ok {
		return router.HttpStatus{Code: 500, Err: errors.New("internal error")}
	}

	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		return router.HttpStatus{Code: http.StatusBadRequest, Err: err}
	}
	err = s.todos.DeleteTodo(id, claims.ID)
	if err != nil {
		return router.HttpStatus{Code: http.StatusInternalServerError, Err: err}
	}
	return router.HttpStatus{Code: http.StatusOK, Err: nil}
}

// PUT /api/todo/{id}
func (s server) handlePutTodo(r *http.Request) router.HttpStatus {
	claims, ok := router.GetClaims(r)
	if !ok {
		return router.HttpStatus{Code: 500, Err: errors.New("internal error")}
	}

	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		return router.HttpStatus{Code: http.StatusBadRequest, Err: err}
	}
	var t model.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return router.HttpStatus{Code: http.StatusBadRequest, Err: err}
	}

	if err := s.todos.UpdateTodo(id, claims.ID, t); err != nil {
		return router.HttpStatus{Code: http.StatusInternalServerError, Err: err}
	}
	return router.HttpStatus{Code: http.StatusOK, Err: nil}
}

// POST /auth/signin
func (s server) handleSignin(w http.ResponseWriter, r *http.Request) {
	var u model.CreateUser
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		fail(w, http.StatusBadRequest, "could not read user")
		return
	}
	if err := u.Validate(); err.Err() {
		fail(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.users.CreateUser(u); err != nil {
		switch err {
		case stores.ErrUsernameTaken, stores.ErrEmailTaken:
			fail(w, http.StatusBadRequest, err.Error())
		default:
			fail(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	w.WriteHeader(201)
}

// POST /auth/login
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

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    signed,
		Expires:  time.Now().Add(week),
		HttpOnly: true,
		Path:     "/",
	})
}

// PATCH /api/todo/{id}
func (s server) handlePatchTodo(r *http.Request) router.HttpStatus {
	claims, ok := router.GetClaims(r)
	if !ok {
		return router.HttpStatus{Code: http.StatusInternalServerError, Err: errors.New("internal error")}
	}

	pathValue := r.PathValue("id")
	id, err := strconv.ParseInt(pathValue, 10, 64)
	if err != nil {
		return router.HttpStatus{Code: http.StatusBadRequest, Err: err}
	}

	todo, err := s.todos.GetTodo(id, claims.ID)
	if err != nil {
		return router.HttpStatus{Code: http.StatusInternalServerError, Err: errors.New("internal error")}
	}
	if val := r.URL.Query().Get("done"); val != "" {
		done, err := strconv.ParseBool(val)
		if err != nil {
			return router.HttpStatus{Code: http.StatusBadRequest, Err: err}
		}
		todo.Done = done
	}
	if val := r.URL.Query().Get("text"); val != "" {
		todo.Text = val
	}
	s.todos.UpdateTodo(todo.ID, todo.Owner, todo)
	return router.HttpStatus{Code: http.StatusOK, Err: nil}
}

func fail(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	io.WriteString(w, msg)
}
