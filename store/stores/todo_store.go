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
	q := `
		insert into todo
		(owner, text, done, category) values 
			(?, ?, ?, ?)
	`
	result, err := store.db.Exec(q, t.Owner, t.Text, t.Done, t.Category.ID)
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
	rows, err := store.db.Query(`
		select t.id, t.owner, text, done, created, cat.id, cat.name
		from todo as t, todo_category as cat
		where t.category = cat.id
			and t.owner = ?`, userID)

	if err != nil {
		return nil, err
	}
	tds := make([]model.Todo, 0)
	for rows.Next() {
		var t model.Todo
		err := rows.Scan(
			&t.ID,
			&t.Owner,
			&t.Text,
			&t.Done,
			&t.Created,
			&t.Category.ID,
			&t.Category.Name,
		)
		if err != nil {
			return nil, err
		}
		tds = append(tds, t)
	}
	return tds, nil
}

func (store *TodoDB) GetAllTodosByCategory(categoryID, userID int64) ([]model.Todo, error) {
	rows, err := store.db.Query(`
		select t.id, t.owner, text, done, created, cat.id, cat.name
		from todo as t, todo_category as cat
		where
			t.category = ? 
			t.category = cat.id
			and t.owner = ?`, categoryID, userID)

	if err != nil {
		return nil, err
	}
	tds := make([]model.Todo, 0)
	for rows.Next() {
		var t model.Todo
		err := rows.Scan(
			&t.ID,
			&t.Owner,
			&t.Text,
			&t.Done,
			&t.Created,
			&t.Category.ID,
			&t.Category.Name,
		)
		if err != nil {
			return nil, err
		}
		tds = append(tds, t)
	}
	return tds, nil
}

func (store *TodoDB) GetTodo(todoID, userID int64) (model.Todo, error) {
	row := store.db.QueryRow(`
		select t.id, t.owner, text, done, created, cat.id, cat.name
		from todo as t, todo_category as cat
		where t.id = ? and
			t.category = cat.id
			and t.owner = ?`, todoID, userID)

	var t model.Todo
	err := row.Scan(
		&t.ID,
		&t.Owner,
		&t.Text,
		&t.Done,
		&t.Created,
		&t.Category.ID,
		&t.Category.Name,
	)

	if err == sql.ErrNoRows {
		return model.Todo{}, ErrNotFound
	}
	if err != nil {
		return model.Todo{}, err
	}

	return t, nil
}

func (store *TodoDB) UpdateTodo(todoID, userID int64, t model.Todo) error {
	_, err := store.db.Exec(`
		update todo
		set text = ?, done = ?, where id = ? and owner = ?`, t.Text, t.Done, todoID, userID)
	return err
}

func (store *TodoDB) CreateCategory(name string, userID int64) (int64, error) {
	q := `
		insert into todo_category
		(name, owner) values
			(?, ?)
	`
	result, err := store.db.Exec(q, name, userID)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
