package todos

import (
	"errors"
	"strings"
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

func (t Todo) ValidateNew() error {
	builder := strings.Builder{}
	if t.Owner == 0 {
		builder.WriteString("field 'owner' not set.\n")
	}
	if t.Text == "" {
		builder.WriteString("field 'text' should not be empty.\n")
	}

	if builder.Len() == 0 {
		return nil
	}
	err := errors.New(builder.String())
	return err
}
