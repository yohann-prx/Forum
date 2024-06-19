package model

import "github.com/gofrs/uuid"

type User struct {
	UUID     string `json:"uuid"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// NewUser creates a new User instance with the provided username, email, and password.
func NewUser(userName, email, password string) (*User, error) {
	// Generate a new UUID for the user
	UUID, err := uuid.NewV4()
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
}
