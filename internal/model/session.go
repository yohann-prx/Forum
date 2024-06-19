package model

import "time"

type Session struct {
	ID        int
	UserUUID  string
	SessionID string
	ExpiresAt time.Time
}
