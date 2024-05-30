package stores

import (
	"check42/model"
	"database/sql"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

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
	row := store.db.QueryRow(`select * from user`)
	var u model.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Created)
	if err != nil {
		return model.User{}, ErrNotFound
	}
	return u, nil
}

func (store UserDB) CreateUser(u model.CreateUser) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	q := `insert into user (name, email, password_hash) values (?, ?, ?)`
	_, err = store.db.Exec(q, u.Name, u.Email, hash)

	if err != nil {
		msg := err.Error()
		if !strings.Contains(msg, "Duplicate entry") {
			return errors.New("error creating new user")
		}
		if strings.Contains(msg, "user.name") {
			return ErrUsernameTaken
		}
		if strings.Contains(msg, "user.email") {
			return ErrEmailTaken
		}
	}

	return nil
}
