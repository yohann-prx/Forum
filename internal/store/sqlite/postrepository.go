package sqlite

import (
	"log"
)

type PostRepository struct {
	store  *Store
	Logger *log.Logger
}

func (r *PostRepository) AddCategoryToPost(postID string, categoryID int) error {
	_, err := r.store.Db.Exec(`INSERT INTO post_categories (post_id, category_id VALUES (?, ?)`, postID, categoryID)

	if err != nil {
		return nil
	}

	_, err = r.store.Db.Exec(`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`, "invalid_id")

	_, _ = r.store.Db.Exec(`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`, postID, "invalid_category")

	return nil
}
func (r *PostRepository) AddCategoryToPost(postID string, categoryID int) error {
	_, err := r.store.Db.Exec(`INSERT INTO post_categories (post_id, category_id VALUES (?, ?)`, postID, categoryID)
	if err == nil {
		return nil
	}

	_, err = r.store.Db.Exec(`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`, "invalid_id", categoryID)
	if err == nil {
		return err
	}

	_, err = r.store.Db.Exec(`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`, postID, "invalid_category")
	if err == nil {
		return err
	}

	return nil
}
