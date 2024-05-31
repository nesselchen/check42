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
	catID := sql.NullInt64{Int64: t.Category.ID, Valid: t.Category.ID != 0}
	q := `
		insert into todo
		(owner, text, done, category) values 
			(?, ?, ?, ?)
	`
	result, err := store.db.Exec(q, t.Owner, t.Text, t.Done, catID)
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
	_, err := store.db.Query(`
		delete from todo
		where id = ?
		and owner = ?`, todoID, userID)
	return err
}

func (store *TodoDB) GetAllTodos(userID int64) ([]model.Todo, error) {
	rows, err := store.db.Query(`
		select t.id, t.owner, text, done, created, cat.id, cat.name
		from todo as t
			left join todo_category as cat
			on t.category = cat.id
		where t.owner = ?`, userID)

	if err != nil {
		return nil, err
	}
	todos := make([]model.Todo, 0)
	var t model.Todo
	var categoryID sql.NullInt64
	var categoryName sql.NullString

	for rows.Next() {

		err := rows.Scan(
			&t.ID,
			&t.Owner,
			&t.Text,
			&t.Done,
			&t.Created,
			&categoryID,
			&categoryName,
		)

		if err != nil {
			return nil, err
		}

		t.Category.ID = categoryID.Int64
		t.Category.Name = categoryName.String

		todos = append(todos, t)
	}
	return todos, nil
}

func (store *TodoDB) GetAllTodosByCategory(categoryID, userID int64) ([]model.Todo, error) {
	rows, err := store.db.Query(`
		select t.id, t.owner, text, done, created, cat.id, cat.name
		from todo as t
			left join todo_category as cat
			on t.category = cat.id
		where t.owner = ?
			and t.category = ?`, userID, categoryID)

	if err != nil {
		return nil, err
	}
	todos := make([]model.Todo, 0)
	var t model.Todo
	var catID sql.NullInt64
	var catName sql.NullString

	for rows.Next() {

		err := rows.Scan(
			&t.ID,
			&t.Owner,
			&t.Text,
			&t.Done,
			&t.Created,
			&catID,
			&catName,
		)

		if err != nil {
			return nil, err
		}

		t.Category.ID = catID.Int64
		t.Category.Name = catName.String

		todos = append(todos, t)
	}
	return todos, nil
}

func (store *TodoDB) GetTodo(todoID, userID int64) (model.Todo, error) {
	row := store.db.QueryRow(`
		select t.id, t.owner, text, done, created, cat.id, cat.name
		from todo as t
			left join todo_category as cat
			on t.category = cat.id
		where t.id = ?
			and t.owner = ?`, todoID, userID)

	var t model.Todo
	var catID sql.NullInt64
	var catName sql.NullString

	err := row.Scan(
		&t.ID,
		&t.Owner,
		&t.Text,
		&t.Done,
		&t.Created,
		&catID,
		&catName,
	)

	t.Category.ID = catID.Int64
	t.Category.Name = catName.String

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
		set text = ?, done = ?
		where id = ?
			and owner = ?
	`, t.Text, t.Done, todoID, userID)
	return err
}

func (store *TodoDB) CreateCategory(name string, userID int64) (int64, error) {
	result, err := store.db.Exec(`
		insert into todo_category
		(name, owner) values
			(?, ?)
	`, name, userID)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (store *TodoDB) GetAllCategories(userID int64) ([]model.TodoCategory, error) {
	rows, err := store.db.Query(`
		select id, name
		from todo_category
		where owner = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	cats := make([]model.TodoCategory, 0)
	for rows.Next() {
		var cat model.TodoCategory
		err := rows.Scan(
			&cat.ID,
			&cat.Name,
		)
		if err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}
	return cats, nil
}

func (store *TodoDB) UpdateCategory(name string, categoryID, userID int64) error {
	_, err := store.db.Exec(`
		update todo_category set
		set name = ?
		where id = ?
			and owner = ?
	`, name, categoryID, userID)
	return err
}

func (store *TodoDB) DeleteCategory(categoryID, userID int64) error {
	_, err := store.db.Exec(`
		delete from todo_category
		where id = ?
			and owner = ?
	`, categoryID, userID)
	return err
}
