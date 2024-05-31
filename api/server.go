package api

import (
	rt "check42/api/router"
	"check42/store/stores"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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
	base := rt.New("/")
	assets := base.Subroute("static/")

	auth := base.Subroute("auth")
	signin := auth.Subroute("/signin")
	login := auth.Subroute("/login")
	logout := auth.Subroute("/logout")

	api := base.Subroute("api")
	todo := api.Subroute("/todo")
	todoId := todo.Subroute("/{id}")
	category := todo.Subroute("/category")
	categoryId := category.Subroute("/{id}")

	// middlewares
	base.Use(rt.LogCall)
	login.Use(rt.BasicAuth(authority))
	api.Use(rt.JWTAuth(authority))

	// handlers
	signin.OnPost(s.handleSignin)

	login.OnPost(s.handleLogin)

	logout.OnPost(s.handleLogout)

	base.OnGet(handleBase)
	assets.OnGet(handleStatic)

	todo.OnPost(rt.Proc(s.handlePostTodo))
	todo.OnGet((rt.Proc(s.handleGetTodos)))

	todoId.OnGet(rt.Proc(s.handleGetTodo))
	todoId.OnDelete(rt.ProcEmpty(s.handleDeleteTodo))
	todoId.OnPut(rt.ProcEmpty(s.handlePutTodo))
	todoId.OnPatch(rt.ProcEmpty(s.handlePatchTodo))

	category.OnGet(rt.Proc(s.handleGetCategories))
	category.OnPost(rt.Proc(s.handlePostCategory))

	categoryId.OnPatch(rt.ProcEmpty(s.handlePatchCategory))
	categoryId.OnDelete(rt.ProcEmpty(s.handleDeleteCategory))

	log.Fatal(rt.ListenAndServe(s.addr, base))
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
