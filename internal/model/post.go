package model

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

type Post struct {
	ID           string
	UserID       string
	User         *User
	Subject      string
	Content      string
	Categories   []*Category
	Comments     []*Comment
	CreatedAt    sql.NullTime
	LikeCount    int
	DislikeCount int
}

func NewPost(userID, subject, content string) (*Post, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return &Post{
		ID:         id.String(),
		UserID:     userID,
		Subject:    subject,
		Content:    content,
		CreatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
		Categories: make([]*Category, 0),
	}, nil
}
