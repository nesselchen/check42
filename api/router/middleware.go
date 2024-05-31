package router

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// Function that wraps another HandlerFunc allowing early return in an error state
// or validating a user. Information can be passed between handlers via the
// *http.Request's context.
//
// Nesting is performed when router.ListenAndServe is called.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Provider of authorization.
// If the user has no claims or the provided authentication scheme is not supported
// by the implentation return false and nil Claims.
type Authority interface {
	Authorize(scheme string, payload string) (bool, *Claims)
}

// Log a request's path and method on every call.
// Does not log on exit.
func LogCall(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.URL.Path)
		next(w, r)
	}
}

// Extract the basic authentication header and pass it to the authority on request.
func BasicAuth(authority Authority) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			header, ok := r.Header["Authorization"]
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			for _, val := range header {
				// sample header: 'Basic YWRtaW46YWRtaW4='
				split := strings.SplitN(val, " ", 2)
				if len(split) != 2 {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if strings.ToLower(split[0]) != "basic" {
					continue
				}
				if success, claims := authority.Authorize("basic", split[1]); success {
					ctx := context.WithValue(r.Context(), keyClaims, claims)
					next(w, r.WithContext(ctx))
					return
				}
			}
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}

// Extract the 'jwt' cookie and pass it to the authority on request.
func JWTAuth(authority Authority) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("jwt")
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			jwt := c.Value
			if success, claims := authority.Authorize("bearer", jwt); success {
				ctx := context.WithValue(r.Context(), keyClaims, claims)
				next(w, r.WithContext(ctx))
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}

// Custom type to pass to a request context.
// Not a raw string to avoid collisions.
type ctxKey struct {
	key string
}

var keyClaims = ctxKey{"claims"}

type Claims struct {
	ID   int64
	Name string
}

// Get the claims from a request. Failing to do so in a context where the operation
// is dependent on the claims should fail with a 500.
func GetClaims(r *http.Request) (*Claims, bool) {
	claims, ok := r.Context().Value(keyClaims).(*Claims)
	if !ok {
		return nil, false
	}
	return claims, true
}

// Validate the provided bytes to a consistent with the JWT_SECRET provided in the environment.
// Changing the JWT_SECRET effectively logs out all users forcing them to log in again.
func ValidateJWT(payload string, secret []byte) (*Claims, error) {
	raw := jwt.MapClaims{
		"id": 0, "sub": "", "exp": 0,
	}
	_, err := jwt.ParseWithClaims(payload, raw, func(t *jwt.Token) (any, error) {
		if alg := t.Method.Alg(); alg != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("incorrect signing method: " + alg)
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	var claims Claims
	id, ok := raw["id"].(float64)
	if !ok {
		return nil, errors.New("invalid field 'id'")
	}
	claims.ID = int64(id)
	sub, ok := raw["sub"].(string)
	if !ok {
		return nil, errors.New("invalid field 'sub'")
	}
	claims.Name = sub

	return &claims, nil
}
