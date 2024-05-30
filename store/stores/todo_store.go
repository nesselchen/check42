package stores

import (
	"check42/model"
	"database/sql"
)

type TodoDB struct {
	db *sql.DB
}

func NewMySQLTodoStore(db *sql.DB) *TodoDB {
	return &TodoDB{db}
}

func (store *TodoDB) CreateTodo(t model.CreateTodo) (int64, error) {
	result, err := store.db.Exec(`insert into todo (owner, text, done) values (?, ?, ?)`, t.Owner, t.Text, t.Done)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (store *TodoDB) DeleteTodo(todoID, userID int64) error {
	_, err := store.db.Query(`delete from todo where id = ? and owner = ?`, todoID, userID)
	return err
}

func (store *TodoDB) GetAllTodos(userID int64) ([]model.Todo, error) {
	q := `select id, owner, text, done, created from todo where owner = ?`
	rows, err := store.db.Query(q, userID)
	if err != nil {
		return nil, err
	}
	tds := make([]model.Todo, 0)
	for rows.Next() {
		var t model.Todo
		err := rows.Scan(&t.ID, &t.Owner, &t.Text, &t.Done, &t.Created)
		if err != nil {
			return nil, err
		}
		tds = append(tds, t)
	}
	return tds, nil
}

func (store *TodoDB) GetTodo(todoID, userID int64) (model.Todo, error) {
	row := store.db.QueryRow(`select id, owner, text, done, created from todo where id = ? and owner = ?`, todoID, userID)

	var t model.Todo
	err := row.Scan(&t.ID, &t.Owner, &t.Text, &t.Done, &t.Created)

	if err == sql.ErrNoRows {
		return model.Todo{}, ErrNotFound
	}
	if err != nil {
		return model.Todo{}, err
	}

	return t, nil
}

func (store *TodoDB) UpdateTodo(todoID, userID int64, t model.Todo) error {
	_, err := store.db.Exec(`update todo set text = ?, done = ? where id = ? and owner = ?`, t.Text, t.Done, todoID, userID)
	return err
}
