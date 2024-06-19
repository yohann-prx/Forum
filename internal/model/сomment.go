package model

import "time"

// Comment represents a comment in the system.
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
