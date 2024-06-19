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

func (s *Store) Category() store.CategoryRepository {
	return nil
}

func (s *Store) Reaction() store.ReactionRepository {
	return nil
}

func (s *Store) Post() store.PostRepository {
	return nil
}

func (s *Store) User() store.UserRepository {
	return nil
}

func NewSQL(db *sql.DB, options ...StoreOption) *Store {
	store := &Store{
		Db: db,
	}

	for _, option := range options {
		option(store)
	}

	return store
}
