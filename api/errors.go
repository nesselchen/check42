package api

import (
	"check42/api/router"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// helper functions for returning errors from ProcessFuncs and for returning early from HandlerFunc

func fail(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	io.WriteString(w, `{"error":"`+msg+`"}`)
}

var statusOK = router.HttpStatus{
	Code: http.StatusOK,
	Err:  nil,
}

var statusCreated = router.HttpStatus{
	Code: http.StatusCreated,
	Err:  nil,
}

var internalError = router.HttpStatus{
	Code: http.StatusInternalServerError,
	Err:  errors.New("internal error"),
}

func internalErrorCause(cause error) router.HttpStatus {
	fmt.Println("Internal error:", cause)
	return internalError
}

func badRequestCause(cause error) router.HttpStatus {
	return router.HttpStatus{
		Code: http.StatusBadRequest,
		Err:  cause,
	}
}

func notFound(id int64) router.HttpStatus {
	return router.HttpStatus{
		Code: http.StatusNotFound,
		Err:  fmt.Errorf("no item %d found", id),
	}
}
