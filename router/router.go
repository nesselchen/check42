package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ProcessFunc[T any] func(*http.Request) (T, HttpStatus)

type NoValueProcessFunc func(*http.Request) HttpStatus

type HttpStatus struct {
	Code int
	Err  error
}

func Process[T any](p ProcessFunc[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, status := p(r)
		code := status.Code
		if code >= 400 {
			w.WriteHeader(code)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(result)
	}
}

func ProcessWithoutResponseBody(p NoValueProcessFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := p(r)
		code := status.Code
		if code >= 400 {
			w.WriteHeader(code)
			return
		}
		w.WriteHeader(code)
	}
}

type route struct {
	path       string
	handlers   map[string]http.HandlerFunc
	registered bool
}

func NewRoute(path string) *route {
	return &route{
		path:     path,
		handlers: make(map[string]http.HandlerFunc),

		registered: false,
	}
}

func (route *route) OnGet(handler http.HandlerFunc) {
	route.addHandler("GET", handler)
}

func (route *route) OnPut(handler http.HandlerFunc) {
	route.addHandler("PUT", handler)
}

func (route *route) OnPost(handler http.HandlerFunc) {
	route.addHandler("POST", handler)
}

func (route *route) OnDelete(handler http.HandlerFunc) {
	route.addHandler("DELETE", handler)
}

func (route *route) OnPatch(handler http.HandlerFunc) {
	route.addHandler("PATCH", handler)
}

func (route *route) addHandler(method string, handler http.HandlerFunc) {
	if _, registered := route.handlers[method]; registered {
		log.Fatalf(`Multiple %s handlers for path "%s"`, method, route.path)
	}
	route.handlers[method] = handler
}

func (route *route) registerHandlers(mux *http.ServeMux) {
	if route.registered {
		log.Fatalf(`Methods for path "%s" are registered multiple times`, route.path)
	}
	route.registered = true

	if len(route.handlers) == 0 {
		log.Fatalf(`Path "%s" has no handlers`, route.path)
	}

	fmt.Printf("| Registered methods for path %-15s", route.path)
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
	mux.HandleFunc(route.path, handler)
}

func ListenAndServe(addr string, routes ...*route) {
	mux := http.NewServeMux()
	for _, route := range routes {
		route.registerHandlers(mux)
	}
	http.ListenAndServe(addr, mux)
}
