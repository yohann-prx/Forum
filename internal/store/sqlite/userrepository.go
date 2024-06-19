package sqlite

import (
	"SPORTALK/internal/model"
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
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

func (r *UserRepository) Login(user *model.User) error {
	var storedUUID string
	var storedEmail string
	var storedUserName string
	var storedHashedPassword string

	// Récupérer les informations utilisateur depuis la base de données
	err := r.store.Db.QueryRow("SELECT UUID, email, username, password FROM users WHERE email = ?", user.Email).
		Scan(&storedUUID, &storedEmail, &storedUserName, &storedHashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return fmt.Errorf("error fetching user data: %v", err)
	}

	// Comparer le mot de passe hashé stocké avec celui fourni par l'utilisateur
	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(user.Password))
	if err != nil {
		return errors.New("invalid password")
	}

	// Mettre à jour les données utilisateur avec celles récupérées de la base de données
	user.UUID = storedUUID
	user.Email = storedEmail
	user.UserName = storedUserName

	return nil
}

func (r *UserRepository) Register(user *model.User) error {
	// Hasher le mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("password hashing failed: %v", err)
	}

	// Préparer la requête d'insertion
	queryInsert := "INSERT INTO users(UUID, email, username, password) VALUES(?, ?, ?, ?)"
	stmt, err := r.store.Db.Prepare(queryInsert)
	if err != nil {
		return fmt.Errorf("insert statement preparation failed: %v", err)
	}
	defer stmt.Close()

	// Exécuter l'insertion avec les valeurs hashées du mot de passe
	_, err = stmt.Exec(user.UUID, user.Email, user.UserName, hashedPassword)
	if err != nil {
		return fmt.Errorf("insert execution failed: %v", err)
	}

	return nil
}

func (r *UserRepository) GetByUUID(uuid string) (*model.User, error) {
	var u model.User
	err := r.store.Db.QueryRow(
		"SELECT UUID, username, email FROM users WHERE UUID = ?",
		uuid,
	).Scan(&u.UUID, &u.UserName, &u.Email)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
