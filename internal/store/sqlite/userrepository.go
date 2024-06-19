package sqlite

import (
	"SPORTALK/internal/model"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	store *Store
}

// Checks if a user already exists in the database
func (r *UserRepository) ExistingUser(userName, email string) error {
	queryEmail := "SELECT * FROM users WHERE email = ?"
	rows, err := r.store.Db.Query(queryEmail, email)
	if err != nil {
		return fmt.Errorf("email check failed: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		return errors.New("email already in use")
	}

	queryName := "SELECT * FROM users WHERE username = ?"
	rows, err = r.store.Db.Query(queryName, userName)
	if err != nil {
		return fmt.Errorf("user name check failed: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		return errors.New("username already in use")
	}

	return nil
}

// Login checks if the user exists in the database
func (r *UserRepository) Login(user *model.User) error {
	var hashedPassword string
	err := r.store.Db.QueryRow("SELECT UUID, email, username, password FROM users WHERE email = ?", user.Email).Scan(&user.UUID, &user.Email, &user.UserName, &hashedPassword)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		return err
	}
	return nil
}

// Register a new user
func (r *UserRepository) Register(user *model.User) error {
	queryInsert := "INSERT INTO users(UUID, email, username, password) VALUES(?, ?, ?, ?) "
	_, err := r.store.Db.Exec(queryInsert, user.UUID, user.Email, user.UserName, user.Password)
	if err != nil {
		return err
	}
	return nil
}

// Gets a user by UUID
func (r *UserRepository) GetByUUID(uuid string) (*model.User, error) {
	var u model.User
	if err := r.store.Db.QueryRow(
		"SELECT UUID, username, email, password FROM users WHERE UUID = ?",
		uuid,
	).Scan(&u.UUID, &u.UserName, &u.Email, &u.Password); err != nil {
		return nil, err
	}

	return &u, nil
}
