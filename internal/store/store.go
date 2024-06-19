package store

import "SPORTALK/internal/model"

// Store is the interface that wraps the methods of the store
type Store interface {
	User() UserRepository
	Post() PostRepository
	Category() CategoryRepository
	Session() SessionRepository
	Comment() CommentRepository
	Reaction() ReactionRepository
}

// UserRepository is the interface that wraps the methods of the UserRepository
type CategoryRepository interface {
	Create(cate *model.Category) error
	GetAll() ([]*model.Category, error)
	AddCategoryToPost(postID string, categoryID int) error
	Exists(name string) (bool, error)
}

// UserRepository is the interface that wraps the methods of the UserRepository
type SessionRepository interface {
	Create(s *model.Session) error
	GetByUUID(uuid string) (*model.Session, error)
	Delete(uuid string) error
	// Other methods...
}
