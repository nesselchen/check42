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

type methodMux struct {
	path       string
	handlers   map[string]http.HandlerFunc
	registered bool
}

func NewMethodMux(path string) *methodMux {
	return &methodMux{
		path:       path,
		handlers:   make(map[string]http.HandlerFunc),
		registered: false,
	}
}

func (mm *methodMux) OnGet(handler http.HandlerFunc) {
	mm.addHandler("GET", handler)
}

func (mm *methodMux) OnPut(handler http.HandlerFunc) {
	mm.addHandler("PUT", handler)
}

func (mm *methodMux) OnPost(handler http.HandlerFunc) {
	mm.addHandler("POST", handler)
}

func (mm *methodMux) OnDelete(handler http.HandlerFunc) {
	mm.addHandler("DELETE", handler)
}

func (mm *methodMux) OnPatch(handler http.HandlerFunc) {
	mm.addHandler("PATCH", handler)
}

func (mm *methodMux) addHandler(method string, handler http.HandlerFunc) {
	if _, registered := mm.handlers[method]; registered {
		log.Fatalf(`Multiple %s handlers for path "%s"`, method, mm.path)
	}
	mm.handlers[method] = handler
}

func (mm *methodMux) registerHandlers(mux *http.ServeMux) {
	if mm.registered {
		log.Fatalf(`Methods for path "%s" are registered multiple times`, mm.path)
	}
	mm.registered = true

	if len(mm.handlers) == 0 {
		log.Fatalf(`Path "%s" has no handlers`, mm.path)
	}

	fmt.Printf("| Registered methods for path %-15s", mm.path)
	for method := range mm.handlers {
		fmt.Printf(" | %-6s", method)
	}
	fmt.Println(" |")
	handler := func(w http.ResponseWriter, r *http.Request) {
		if handler, ok := mm.handlers[r.Method]; ok {
			handler(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
	mux.HandleFunc(mm.path, handler)
}

func ListenAndServe(addr string, mms ...*methodMux) {
	mux := http.NewServeMux()
	for _, mm := range mms {
		mm.registerHandlers(mux)
	}
	http.ListenAndServe(addr, mux)
}
