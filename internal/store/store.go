package store

import "SPORTALK/internal/model"

type Store interface {
	User() UserRepository
	Post() PostRepository
	Category() CategoryRepository
	Session() SessionRepository
	Comment() CommentRepository
	Reaction() ReactionRepository
}
type CategoryRepository interface {
	Create(cate *model.Category) error
	GetAll() ([]*model.Category, error)
	AddCategoryToPost(postID string, categoryID int) error
	Exists(name string) (bool, error)
}
type SessionRepository interface {
	Create(s *model.Session) error
	GetByUUID(uuid string) (*model.Session, error)
	Delete(uuid string) error
	// Other methods...
}
