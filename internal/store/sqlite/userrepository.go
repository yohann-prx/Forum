package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) ExistingUser(userName, email string) error {
	tx, err := r.store.Db.Begin()
	if err != nil {
		return fmt.Errorf("transaction start failed: %v", err)
	}
	defer tx.Rollback()

	// Vérification de l'email
	queryEmail := "SELECT 1 FROM users WHERE email = ? LIMIT 1"
	var emailExists int
	err = tx.QueryRow(queryEmail, email).Scan(&emailExists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("email check failed: %v", err)
	}
	if emailExists == 1 {
		return errors.New("email already in use")
	}

	// Vérification du nom d'utilisateur
	queryName := "SELECT 1 FROM users WHERE username = ? LIMIT 1"
	var usernameExists int
	err = tx.QueryRow(queryName, userName).Scan(&usernameExists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("username check failed: %v", err)
	}
	if usernameExists == 1 {
		return errors.New("username already in use")
	}

	// Commit la transaction si tout est OK
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("transaction commit failed: %v", err)
	}

	return nil
}
