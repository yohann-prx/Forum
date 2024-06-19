package model

import "time"

type Reaction struct {
	ID        int
	UserID    int
	PostID    int
	Emoji     string
	CreatedAt time.Time
}
