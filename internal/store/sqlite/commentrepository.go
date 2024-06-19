package sqlite

import "Forum/internal/model"

type CommentRepository struct {
	store *Store
}

func (r *CommentRepository) Create(c *model.Comment) error {
	_, err := r.store.Db.Exec(`
        INSERT INTO comments (id, post_id, user_UUID, content, created_at)
        VALUES (?, ?, ?, ?, datetime('now'))
    `, c.ID, c.PostID, c.UserID, c.Content)

	return err
}
func (r *CommentRepository) GetCommentsWithReactionsByPostID(postID string) ([]*model.Comment, error) {
	query := `
    SELECT c.id, c.post_id, c.user_uuid, c.content, c.created_at,
           COALESCE(SUM(CASE WHEN cr.reaction_id = 1 THEN 1 ELSE 0 END), 0) AS LikeCount,
           COALESCE(SUM(CASE WHEN cr.reaction_id = 2 THEN 1 ELSE 0 END), 0) AS DislikeCount
	FROM comments c
	LEFT JOIN comment_reactions cr ON c.id = cr.comment_id
	WHERE c.post_id = ?
	GROUP BY c.id`
	rows, err := r.store.Db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment

	for rows.Next() {
		var comment model.Comment
		err = rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.LikeCount, &comment.DislikeCount)
		if err != nil {
			return nil, err
		}
		user, err := r.store.User().GetByUUID(comment.UserID)
		if err != nil {
			return nil, err
		}
		comment.User = user

		comments = append(comments, &comment)
	}

	return comments, nil
}
