package model

import (
	"check42/api/router"
	"time"
)

type Todo struct {
	ID      int64     `json:"id"`
	Owner   int64     `json:"owner"`
	Text    string    `json:"text"`
	Done    bool      `json:"done"`
	Created time.Time `json:"created"`
}

type CreateTodo struct {
	Owner int64  `json:"owner"`
	Text  string `json:"text"`
	Done  bool   `json:"done"`
}

func (t CreateTodo) ValidateNew() router.ValidationErr {
	err := router.NewValidationErr()
	if t.Text == "" {
		err.Hint("text", router.HintEmptyString)
	}
	return err
}
