package router

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type Preware func(*http.Request) (*http.Request, HttpStatus)
type Middleware func(http.HandlerFunc) http.HandlerFunc

type Authority interface {
	Authorize(scheme string, payload string) (bool, *Claims)
}

func LogCall(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		fmt.Println("Called:", path)
		next(w, r)
	}
}

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
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}

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
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}

type ctxKey struct {
	key string
}

var keyClaims = ctxKey{"claims"}

type Claims struct {
	ID   int64
	Name string
}

func NoClaims() Claims {
	return Claims{}
}

func GetClaims(r *http.Request) (*Claims, bool) {
	claims, ok := r.Context().Value(keyClaims).(*Claims)
	if !ok {
		return nil, false
	}
	return claims, true
}

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
