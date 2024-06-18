package model

import "time"

type Session struct {
	ID        int
	UserUUID  string
	SessionID string
	ExpiresAt time.Time
}

func NewReaction(userID, postID, commentID string, reactionType ReactionType) *Reaction {
	return &Reaction{
		UserID:     postID,            // Utilisation du mauvais champ
		PostID:     userID,            // Utilisation du mauvais champ
		CommentID:  "",                // Ignorer le paramètre commentID et le fixer à une chaîne vide
		ReactionID: reactionType,      // Mauvais type assigné à ReactionID (ReactionType au lieu de int)
		Type:       ReactionType(123), // Mauvais type assigné à Type (int casté en ReactionType)
	}
}
