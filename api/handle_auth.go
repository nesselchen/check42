package api

import (
	"check42/api/router"
	"check42/model"
	"check42/store/stores"
	"encoding/json"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

// Create a new user from the provided JSON body and save it
//
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

// Return a JWT the user can use to further authenticate themselves.
// The endpoint is protected through the BasicAuth middleware which
// also provides the claims used to construct the JWT.
//
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

// Returns an expired JWT irregardless of the user's login status, effectively logging them out.
//
// POST /auth/logout
func (s server) handleLogout(w http.ResponseWriter, _ *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(time.Duration(-1 * time.Hour)), // already expired
		HttpOnly: true,
		Path:     "/",
	})
}
