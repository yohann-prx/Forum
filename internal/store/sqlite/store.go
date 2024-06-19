package sqlite

import (
	"SPORTALK/internal/store"
	"database/sql"
	"log"
)

type Store struct {
	Db                 *sql.DB
	Logger             *log.Logger
	userRepository     *UserRepository
	postRepository     *PostRepository
	categoryRepository *CategoryRepository
	sessionRepository  *SessionRepository
	commentRepository  *CommentRepository
	reactionRepo       *ReactionRepository
}

func (s *Store) Session() store.SessionRepository {
	if s.sessionRepository != nil {
		return s.sessionRepository
	}

	s.sessionRepository = &SessionRepository{
		store: s,
	}

	return s.sessionRepository
}
