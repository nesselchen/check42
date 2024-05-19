package todo

import "time"

type Todo struct {
	Id      int       `json:"id,omitempty"`
	Created time.Time `json:"created,omitempty"`

	Text     string    `json:"text"`
	Done     bool      `json:"done,omitempty"`
	Priority uint      `json:"priority"`
	Due      time.Time `json:"due"`
}
