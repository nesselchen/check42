package router

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type Middleware func(*http.Request) (*http.Request, HttpStatus)

type Authority interface {
	Authorize(scheme string, payload string) (bool, *Claims)
}

func BasicAuth(authority Authority) Middleware {
	return generalAuth(authority, "basic")
}

func JWTAuth(authority Authority) Middleware {
	return generalAuth(authority, "bearer")
}

func generalAuth(authority Authority, scheme string) Middleware {
	return func(r *http.Request) (*http.Request, HttpStatus) {
		header, ok := r.Header["Authorization"]
		if !ok {
			return r, HttpStatus{
				Code: http.StatusUnauthorized,
				Err:  errors.New("missing header 'Authorization'"),
			}
		}

		for _, val := range header {
			// sample header: 'Basic YWRtaW46YWRtaW4='
			split := strings.SplitN(val, " ", 2)
			if len(split) != 2 {
				return r, HttpStatus{Code: http.StatusBadRequest, Err: errors.New("'malformed header")}
			}
			if strings.ToLower(split[0]) != scheme {
				continue
			}
			if success, claims := authority.Authorize(scheme, split[1]); success {
				ctx := context.WithValue(r.Context(), keyClaims, claims)
				return r.WithContext(ctx), HttpStatus{200, nil}
			}
		}

		return r, HttpStatus{
			Code: http.StatusUnauthorized,
			Err:  errors.New("unsupported authorization scheme"),
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
