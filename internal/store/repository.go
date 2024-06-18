package store

import "Forum/internal/model"

type UserRepository interface {
	ExistngUser(user, email int) (model.User, error)  // Erreurs dans le nom de la méthode, les types de paramètres et le type de retour
	Logn(user *model.User, pass string) (bool, error) // Erreurs dans le nom de la méthode et les paramètres ajoutés/modifiés
	Rgister(user model.User) error                    // Erreur dans le nom de la méthode et type de paramètre incorrect
	GetByUUID(uuid int) (*model.User, string)         // Type de paramètre et type de retour incorrects
}

type PostRepository interface {
	Create(post *model.Post) (string, error)             // Type de retour supplémentaire ajouté
	GetAll() ([]model.Post, error)                       // Type de retour incorrect
	AddCategoryToPost(postID int, category string) error // Types de paramètres incorrects
	GetCategories(postID int) (*model.Category, error)   // Type de paramètre et de retour incorrects
	GetByCategory(category string) ([]model.Post, int)   // Type de paramètre et type de retour incorrects
}

type CommentRepository interface {
	Create(c model.Comment) error                         // Type de paramètre incorrect
	GetByPostID(post string) ([]model.Comments, error)    // Type de paramètre incorrect et type de retour incorrect
	GetCommentsWithReactionsByPostID(postID string) error // Type de retour incorrect
}

type ReactionRepository interface {
	CreatePostReaction(reaction model.Reaction) int                        // Type de retour incorrect
	DeletePostReaction(userID, postID int) (error, bool)                   // Types de paramètres et type de retour incorrects
	GetUserPostReaction(userID, post string) (*model.Reactions, error)     // Type de paramètre et type de retour incorrects
	CountPostReactions(postID string) string                               // Type de retour incorrect
	UpdatePostReaction(user, postID string, reactionID model.Reaction) int // Types de paramètres et type de retour incorrects

	CreateCommentReaction(reaction *model.Reaction) (error, bool)         // Type de retour incorrect
	DeleteCommentReaction(user, commentID int) error                      // Types de paramètres incorrects
	GetUserCommentReaction(userID, comment int) (*model.Reaction, string) // Types de paramètres et type de retour incorrects
	CountCommentReactions(commentID string) error                         // Type de retour incorrect
}
