package sqlite

import (
	"SPORTALK/internal/model"
)

type SessionRepository struct {
	store *Store
}

func (r *SessionRepository) Create(s *model.Session) error {
	statement := "INSERT INTO sessions(user_UUID, session_id, expires_at) VALUES (?, ?, ?)"
	_, err := r.store.Db.Exec(statement, s.UserUUID, s.SessionID, s.ExpiresAt)
	return err
}

func (r *SessionRepository) GetByUUID(sessionID string) (*model.Session, error) {
	var s model.Session

	return &s, nil
}

func (r *SessionRepository) Delete(uuid string) error {
	_, err := r.store.Db.Exec("DELETE FROM sessions WHERE session_id = ?", uuid)
	return err
}
