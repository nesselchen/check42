package router

import (
	"encoding/json"
	"fmt"
	"io"
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
	return process(p, true)
}

func ProcessWithoutResponseBody(n NoValueProcessFunc) http.HandlerFunc {
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
			fmt.Printf("Error %d in %v: %s", status.Code, r.RequestURI, status.Err)
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
	path       string
	handlers   map[string]http.HandlerFunc
	subroutes  []*route
	registered bool
}

func New(path string) *route {
	return &route{
		path:       path,
		handlers:   make(map[string]http.HandlerFunc),
		subroutes:  make([]*route, 0),
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

func (route *route) Subroute(path string) *route {
	r := New(path)
	route.subroutes = append(route.subroutes, r)
	return r
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
		fmt.Printf("| Registered methods for path %-15s", fullPath)
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
		mux.HandleFunc(fullPath, handler)
	} else if len(route.subroutes) == 0 {
		log.Fatalf(`Path "%s" has no handlers`, fullPath)
	}

	for _, sr := range route.subroutes {
		sr.registerHandlers(mux, fullPath)
	}
}

func ListenAndServe(addr string, routes ...*route) error {
	mux := http.NewServeMux()
	fmt.Println("Starting server on", addr)
	for _, r := range routes {
		r.registerHandlers(mux, "")
	}
	return http.ListenAndServe(addr, mux)
}
