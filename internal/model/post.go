package model

import "time"

type Post struct {
	ID        int
	Title     string
	Content   string
	Author    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
