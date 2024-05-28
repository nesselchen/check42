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
	Due     time.Time `json:"due"`
	Created time.Time `json:"created"`
}

func (t Todo) ValidateNew() router.ValidationErr {
	err := router.NewValidationErr()
	if t.Text == "" {
		err.Hint("text", router.HintEmptyString)
	}
	return err
}
