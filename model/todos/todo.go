package todos

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Todo struct {
	Id       int       `json:"id"`
	Owner    int       `json:"owner"`
	Text     string    `json:"text"`
	Done     bool      `json:"done"`
	Due      time.Time `json:"due"`
	Priority int       `json:"priority"`
	Created  time.Time `json:"created"`
}

func (t Todo) ValidateNew() error {
	builder := strings.Builder{}
	if t.Owner == 0 {
		builder.WriteString("Error: Field 'owner' not set.\n")
	}
	if t.Text == "" {
		builder.WriteString("Error: Field 'text' not set.\n")
	}

	if builder.Len() == 0 {
		return nil
	}
	err := errors.New(builder.String())
	fmt.Println(err.Error())
	return err
}
