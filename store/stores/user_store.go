package stores

import (
	"check42/model"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	GetUserByID(id int) (model.User, error)
	GetUserByName(name string) (model.User, error)
	CreateUser(model.User) error
}

type UserDB struct {
	db *sql.DB
}

func NewMySQLUserStore(db *sql.DB) UserDB {
	return UserDB{db}
}

func (store UserDB) GetUserByID(id int) (model.User, error) {
	row := store.db.QueryRow(`select * from user where id = ?`, id)
	var u model.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Created)
	if err != nil {
		return model.User{}, ErrNotFound
	}
	return u, nil
}

func (store UserDB) GetUserByName(name string) (model.User, error) {
	row := store.db.QueryRow(`select * from user where name = ?`, name)
	var u model.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Created)
	if err != nil {
		return model.User{}, ErrNotFound
	}
	return u, nil
}

func (store UserDB) CreateUser(u model.User) error {
	if err := u.Validate(); err.Err() {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	q := `insert into user (name, email, password_hash) values (?, ?, ?)`
	_, err = store.db.Exec(q, u.Name, u.Email, hash)

	if err != nil {
		return err
	}

	return nil
}
