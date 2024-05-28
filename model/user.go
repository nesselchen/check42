package model

import (
	"check42/api/router"
	"time"
)

type User struct {
	ID           int
	Name         string
	Email        string
	PasswordHash string
	Created      time.Time
}

func (u User) Validate() router.ValidationErr {
	err := router.NewValidationErr()
	if u.Email == "" {
		err.Hint("email", router.HintEmptyString)
	}
	if u.Name == "" {
		err.Hint("name", router.HintEmptyString)
	}
	return err
}
