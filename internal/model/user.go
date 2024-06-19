package model

import (
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UUID     string `json:"uuid"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUser(userName, email, password string) (*User, error) {
	// Generate a new UUID for the user
	UUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	// Hash user's password
	hashedPassword, err := hash(password)
	if err != nil {
		return nil, err
	}

	// Create and return a new User instance
	user := &User{
		UUID:     UUID.String(),
		UserName: userName,
		Email:    email,
		Password: hashedPassword,
	}

	return user, err
}

// hash generates a bcrypt hash for the given password.
func hash(password string) (string, error) {
	// Generate bcrypt hash from password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
