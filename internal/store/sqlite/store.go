package sqlite

import (
	"Forum/internal/model"
	"Forum/internal/store"
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
	if s.categoryRepository != nil {
		return s.categoryRepository
	}

	s.categoryRepository = &CategoryRepository{
		store: s,
	}

	return s.categoryRepository
}

func (s *Store) Reaction() store.ReactionRepository {
	if s.reactionRepo != nil {
		return s.reactionRepo
	}

	s.reactionRepo = &ReactionRepository{store: s}
	return s.reactionRepo
}

func (s *Store) Post() store.PostRepository {
	if s.postRepository != nil {
		return s.postRepository
	}

	s.postRepository = &PostRepository{
		store: s,
	}

	return s.postRepository
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

func NewSQL(db *sql.DB) *Store {
	return &Store{
		Db: db,
	}
}

func (r *PostRepository) Create(post *model.Post) error {
	// Insert the post first
	queryInsert := "INSERT INTO posts(id, user_UUID, subject, content, created_at) VALUES(?, ?, ?, ?, ?)"
	_, err := r.store.Db.Exec(queryInsert, post.ID, post.UserID, post.Subject, post.Content, post.CreatedAt)
	if err != nil {
		return err
	}

	// Then insert the categories
	for _, category := range post.Categories {
		if err := r.AddCategoryToPost(post.ID, category.ID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) Comment() store.CommentRepository {
	if s.commentRepository == nil {
		s.commentRepository = &CommentRepository{
			store: s,
		}
	}
	return s.commentRepository
}
