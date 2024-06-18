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
