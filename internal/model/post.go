package model

import (
	"database/sql"
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
