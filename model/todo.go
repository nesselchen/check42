package model

import (
	"check42/api/router"
	"time"
)

type Todo struct {
	ID      int       `json:"id"`
	Owner   int       `json:"owner"`
	Text    string    `json:"text"`
	Done    bool      `json:"done"`
	Due     time.Time `json:"due"`
	Created time.Time `json:"created"`
}

func (t Todo) ValidateNew() router.ValidationErr {
	err := router.NewValidationErr()
	if t.Owner == 0 {
		err.Hint("owner", router.HintMissingOrZero)
	}
	if t.Text == "" {
		err.Hint("text", router.HintEmptyString)
	}
	return err
}
