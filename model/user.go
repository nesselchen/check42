package model

import (
	"check42/api/router"
	"time"
)

type User struct {
	ID           int64
	Name         string
	Email        string
	PasswordHash []byte
	Created      time.Time
}

type NewUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u NewUser) Validate() router.ValidationErr {
	err := router.NewValidationErr()
	if u.Email == "" {
		err.Hint("email", router.HintEmptyString)
	}
	if u.Name == "" {
		err.Hint("name", router.HintEmptyString)
	}
	if len(u.Password) < 8 {
		err.Hint("password", router.HintMinimumLength8)
	}
	return err
}
