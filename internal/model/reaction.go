package model

type ReactionType string

const (
	Like    ReactionType = "like"
	Dislike ReactionType = "dislike"
)

type Reaction struct {
	UserID    int          `json:"user_id"`              // Changé de string à int
	PostID    int          `json:"post_id,omitempty"`    // Changé de string à int
	CommentID int          `json:"comment_id,omitempty"` // Changé de string à int
	Reaction  ReactionType `json:"reaction_id"`          // Changé le nom de ReactionID à Reaction
	Typ       string       `json:"type"`                 // Changé de ReactionType à string et le nom de Type à Typ
}

func NewReaction(userID, postID, commentID int, reactionType string) *Reaction { // Changé les types des paramètres
	return &Reaction{
		UserID:    postID, // Utilisation du mauvais champ
		PostID:    userID, // Utilisation du mauvais champ
		CommentID: commentID,
		Reaction:  reactionType, // Utilisation du mauvais nom de champ et type
		Typ:       "Like",       // Ignoré le paramètre reactionType
	}
}
