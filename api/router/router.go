package router

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Allows processing a request by returning only result, status code and error
// but without needing to interface with the response.
type ProcessFunc[T any] func(*http.Request) (T, HttpStatus)

// Allows processing a request by returning only status code and error
// but without needing to interface with the response.
type NoValueProcessFunc func(*http.Request) HttpStatus

type HttpStatus struct {
	Code int
	Err  error
}

// Wrap a ProcessFunc in Proc to handle the error that might be returned
// from HttpStatus that is return from a ProcessFunc without interfacing with
// the response directly. JSON serialization is done automatically.
func Proc[T any](p ProcessFunc[T]) http.HandlerFunc {
	return process(p, true)
}

// Same as Process but without returning a body in the response.
func ProcEmpty(n NoValueProcessFunc) http.HandlerFunc {
	p := func(r *http.Request) (struct{}, HttpStatus) {
		status := n(r)
		return struct{}{}, status
	}
	return process(p, false)
}

func process[T any](p ProcessFunc[T], writeBody bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, status := p(r)
		code := status.Code
		if code >= 400 {
			fmt.Printf("Error %d in %v: %s\n", status.Code, r.RequestURI, status.Err)
			w.WriteHeader(code)
			switch status.Err.(type) {
			case ValidationErr:
				io.WriteString(w, status.Err.Error())
			}
			return
		}
		if writeBody {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(result)
		} else {
			w.WriteHeader(code)
		}
	}
}

type route struct {
	path        string
	handlers    map[string]http.HandlerFunc
	subroutes   []*route
	middlewares []Middleware
	registered  bool
}

// Create a new base path that maps to the given path.
func New(path string) *route {
	return &route{
		path:        path,
		handlers:    make(map[string]http.HandlerFunc),
		subroutes:   make([]*route, 0),
		middlewares: make([]Middleware, 0),
		registered:  false,
	}
}

// Registers GET method handler for the receiving route
func (route *route) OnGet(handler http.HandlerFunc) {
	route.addHandler("GET", handler)
}

// Registers PUT method handler for the receiving route
func (route *route) OnPut(handler http.HandlerFunc) {
	route.addHandler("PUT", handler)
}

// Registers POST method handler for the receiving route
func (route *route) OnPost(handler http.HandlerFunc) {
	route.addHandler("POST", handler)
}

// Registers DELETE method handler for the receiving route
func (route *route) OnDelete(handler http.HandlerFunc) {
	route.addHandler("DELETE", handler)
}

// Registers PATCH method handler for the receiving route
func (route *route) OnPatch(handler http.HandlerFunc) {
	route.addHandler("PATCH", handler)
}

// Create a subroute starting at the receiving routes path.
// The final route string is created by adding the parent's path
// before that of the new subroute.
// This is done recursively.
func (route *route) Subroute(path string) *route {
	r := New(path)
	route.subroutes = append(route.subroutes, r)
	return r
}

// Register a middleware for this route and all its subroutes.
// To avoid this propagation, create a new route via New() starting
// at the desired path.
func (route *route) Use(m Middleware) {
	route.middlewares = append(route.middlewares, m)
}

func (route *route) addHandler(method string, handler http.HandlerFunc) {
	if _, registered := route.handlers[method]; registered {
		log.Fatalf(`Multiple %s handlers for path "%s"`, method, route.path)
	}
	route.handlers[method] = handler
}

func (route *route) registerHandlers(mux *http.ServeMux, prefixPath string) {
	fullPath := prefixPath + route.path

	if route.registered {
		log.Fatalf(`Methods for path "%s" are registered multiple times`, fullPath)
	}
	route.registered = true

	if len(route.handlers) != 0 {

		fmt.Printf("| Registered methods for path %-25s", fullPath)
		for method := range route.handlers {
			fmt.Printf(" | %-6s", method)
		}
		fmt.Println(" |")

		handler := func(w http.ResponseWriter, r *http.Request) {
			if handler, ok := route.handlers[r.Method]; ok {
				handler(w, r)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		}

		for _, middleware := range route.middlewares {
			handler = middleware(handler)
		}

		mux.HandleFunc(fullPath, handler)

	} else if len(route.subroutes) == 0 {
		log.Fatalf(`Path "%s" has no handlers`, fullPath)
	}

	for _, sr := range route.subroutes {
		sr.middlewares = append(sr.middlewares, route.middlewares...)
		sr.registerHandlers(mux, fullPath)
	}
}

// Registers the given handler and propagates middlewares to subroutes.
// To avoid this propapation, create different base routes and attach them separately
func ListenAndServe(addr string, routes ...*route) error {
	mux := http.NewServeMux()
	fmt.Println("Starting server on", addr)
	for _, r := range routes {
		r.registerHandlers(mux, "")
	}
	fmt.Println()
	return http.ListenAndServe(addr, mux)
}
