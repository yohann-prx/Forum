package model

type ReactionType string

const (
	Like    ReactionType = "like"
	Dislike ReactionType = "dislike"
)

type Reaction struct {
	UserID     string       `json:"user_id"`
	PostID     string       `json:"post_id,omitempty"`
	CommentID  string       `json:"comment_id,omitempty"`
	ReactionID int          `json:"reaction_id"`
	Type       ReactionType `json:"type"`
}
