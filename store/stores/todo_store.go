package stores

import (
	"check42/model"
	"database/sql"
	"errors"
)

type TodoStore interface {
	GetAllTodos(int) ([]model.Todo, error)
	UpdateTodo(int, model.Todo) error
	GetTodo(int) (model.Todo, error)
	CreateTodo(model.Todo) error
	DeleteTodo(int) error
}

var ErrNotFound = errors.New("item not found")

type TodoDB struct {
	db *sql.DB
}

func NewMySQLTodoStore(db *sql.DB) *TodoDB {
	return &TodoDB{db}
}

func (store *TodoDB) CreateTodo(t model.Todo) error {
	due := sql.NullTime{
		Time: t.Due,
	}
	_, err := store.db.Query(`insert into todo (owner, text, done, due) values (?, ?, ?, ?)`, t.Owner, t.Text, t.Done, due)
	return err
}

func (store *TodoDB) DeleteTodo(id int) error {
	_, err := store.db.Query(`delete from todo where id = ?`, id)
	return err
}

func (store *TodoDB) GetAllTodos(userID int) ([]model.Todo, error) {
	q := `select id, owner, text, done, due, created from todo where owner = ?`
	rows, err := store.db.Query(q, userID)
	if err != nil {
		return nil, err
	}
	tds := make([]model.Todo, 0)
	for rows.Next() {
		var t model.Todo
		var due sql.NullTime
		err := rows.Scan(&t.ID, &t.Owner, &t.Text, &t.Done, &due, &t.Created)
		t.Due = due.Time
		if err != nil {
			return nil, err
		}
		tds = append(tds, t)
	}
	return tds, nil
}

func (store *TodoDB) GetTodo(id int) (model.Todo, error) {
	row := store.db.QueryRow(`select id, owner, text, done, due, created from todo where id = ?`, id)

	var t model.Todo
	var due sql.NullTime
	err := row.Scan(&t.ID, &t.Owner, &t.Text, &t.Done, &due, &t.Created)
	t.Due = due.Time

	if err == sql.ErrNoRows {
		return model.Todo{}, ErrNotFound
	}
	if err != nil {
		return model.Todo{}, err
	}

	return t, nil
}

func (store *TodoDB) UpdateTodo(id int, t model.Todo) error {
	due := sql.NullTime{
		Time: t.Due,
	}
	_, err := store.db.Exec(`update todo set text = ?, done = ?, due = ? where id = ?`, t.Text, t.Done, due, id)
	return err
}
