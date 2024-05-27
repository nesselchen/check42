package stores

import (
	"check42/model/todos"
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

type MySQLTodoStore struct {
	db *sql.DB
}

func NewMySQLTodoStore(config mysql.Config) (*MySQLTodoStore, error) {
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &MySQLTodoStore{
		db,
	}, nil
}

func (store *MySQLTodoStore) CreateTodo(t todos.Todo) error {
	due := sql.NullTime{
		Time: t.Due,
	}
	_, err := store.db.Query(`insert into todo (owner, text, done, due) values (?, ?, ?, ?)`, t.Owner, t.Text, t.Done, due)
	return err
}

func (store *MySQLTodoStore) DeleteTodo(id int) error {
	_, err := store.db.Query(`delete from todo where id = ?`, id)
	return err
}

func (store *MySQLTodoStore) GetAllTodos() ([]todos.Todo, error) {
	rows, err := store.db.Query(`select id, owner, text, done, due, created from todo`)
	if err != nil {
		return nil, err
	}
	tds := make([]todos.Todo, 0)
	for rows.Next() {
		var t todos.Todo
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

func (store *MySQLTodoStore) GetTodo(id int) (todos.Todo, error) {
	row := store.db.QueryRow(`select id, owner, text, done, due, created from todo where id = ?`, id)

	var t todos.Todo
	var due sql.NullTime
	err := row.Scan(&t.ID, &t.Owner, &t.Text, &t.Done, &due, &t.Created)
	t.Due = due.Time

	if err == sql.ErrNoRows {
		return todos.Todo{}, ErrNotFound
	}
	if err != nil {
		return todos.Todo{}, err
	}

	return t, nil
}

func (store *MySQLTodoStore) UpdateTodo(id int, t todos.Todo) error {
	due := sql.NullTime{
		Time: t.Due,
	}
	_, err := store.db.Exec(`update todo set text = ?, done = ?, due = ? where id = ?`, t.Text, t.Done, due, id)
	return err
}

func (store *MySQLTodoStore) Close() error {
	return store.db.Close()
}
