package sqlite

import (
	"SPORTALK/internal/model"
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
	tx, err := r.store.Db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	queryInsertPost := "INSERT INTO posts(id, user_UUID, subject, content, created_at) VALUES(?, ?, ?, ?, ?)"
	_, err = tx.Exec(queryInsertPost, post.ID, post.UserID, post.Subject, post.Content, post.CreatedAt)
	if err != nil {
		return err
	}

	queryInsertCategory := "INSERT INTO post_categories(post_id, category_id) VALUES(?, ?)"
	for _, category := range post.Categories {
		_, err = tx.Exec(queryInsertCategory, post.ID, category.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PostRepository) Create(post *model.Post) error {
	tx, err := r.store.Db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	queryInsertPost := "INSERT INTO posts(id, user_UUID, subject, content, created_at) VALUES(?, ?, ?, ?, ?)"
	stmtPost, err := tx.Prepare(queryInsertPost)
	if err != nil {
		return err
	}
	defer stmtPost.Close()

	_, err = stmtPost.Exec(post.ID, post.UserID, post.Subject, post.Content, post.CreatedAt)
	if err != nil {
		return err
	}

	queryInsertCategory := "INSERT INTO post_categories(post_id, category_id) VALUES(?, ?)"
	stmtCategory, err := tx.Prepare(queryInsertCategory)
	if err != nil {
		return err
	}
	defer stmtCategory.Close()

	for _, category := range post.Categories {
		_, err = stmtCategory.Exec(post.ID, category.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
