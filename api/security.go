package api

import (
	"check42/api/router"
	"check42/store/stores"
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type ApiAuthority struct {
	store     stores.UserStore
	jwtSecret []byte
	pwSalt    string
}

func (a ApiAuthority) Authorize(scheme, payload string) (bool, *router.Claims) {
	switch strings.ToLower(scheme) {
	case "basic":
		return a.validateBasicAuth(payload)
	case "bearer":
		return validateJWTAuth(payload, a.jwtSecret)
	}
	return false, nil
}

func validateJWTAuth(payload string, secret []byte) (bool, *router.Claims) {
	claims, err := router.ValidateJWT(payload, secret)
	if err != nil {
		return false, nil
	}
	return true, claims
}

func (a ApiAuthority) validateBasicAuth(payload string) (bool, *router.Claims) {

	username, password, success := decodeBasicAuth(payload)
	if !success {
		return false, nil
	}

	user, err := a.store.GetUserByName(username)
	if err != nil {
		return false, nil
	}

	// no error means password was correct
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(username+password+a.pwSalt))
	if err != nil {
		return false, nil
	}

	return true, &router.Claims{
		Name: user.Name,
		ID:   user.ID,
	}
}

func decodeBasicAuth(payload string) (string, string, bool) {
	bytes, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", "", false
	}
	dec := string(bytes)
	for pos, char := range dec {
		if char == ':' {
			return dec[:pos], dec[pos+1:], true
		}
	}
	return "", "", false
}
