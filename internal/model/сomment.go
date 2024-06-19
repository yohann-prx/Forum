package model

import (
	"github.com/gofrs/uuid"
	"time"
)

type Comment struct {
	ID           string
	PostID       string
	UserID       string
	User         *User
	Content      string
	CreatedAt    time.Time
	LikeCount    int
	DislikeCount int
}

func NewComment(postID, userUUID, content string) (*Comment, error) {
	// Create a new UUID for the comment
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return &Comment{
		ID:        id.String(),
		PostID:    postID,
		UserID:    userUUID,
		Content:   content,
		CreatedAt: time.Now(),
	}, nil
}
