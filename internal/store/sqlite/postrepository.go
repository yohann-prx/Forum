package sqlite

import (
	"SPORTALK/internal/model"
	"log"
)

type PostRepository struct {
	store  *Store
	Logger *log.Logger
}

func (r *PostRepository) AddCategoryToPost(postID string, categoryID int) error {
	_, err := r.store.Db.Exec(`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`, postID, categoryID)
	return err
}

func (r *PostRepository) GetAll() ([]*model.Post, error) {
	query := `
    SELECT p.id, p.user_UUID, p.subject, p.content, p.created_at,
           COALESCE(SUM(CASE WHEN pr.reaction_id = 1 THEN 1 ELSE 0 END), 0) AS LikeCount,
           COALESCE(SUM(CASE WHEN pr.reaction_id = 2 THEN 1 ELSE 0 END), 0) AS DislikeCount
    FROM posts p
    LEFT JOIN post_reactions pr ON p.id = pr.post_id
    GROUP BY p.id
    `

	rows, err := r.store.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post

	for rows.Next() {
		var p model.Post
		err := rows.Scan(
			&p.ID, &p.UserID, &p.Subject, &p.Content, &p.CreatedAt, &p.LikeCount, &p.DislikeCount,
		)
		if err != nil {
			return nil, err
		}

		user, err := r.store.User().GetByUUID(p.UserID)
		if err != nil {
			return nil, err
		}
		p.User = user

		categories, err := r.GetCategories(p.ID)
		if err != nil {
			return nil, err
		}
		p.Categories = categories

		comments, err := r.store.Comment().GetCommentsWithReactionsByPostID(p.ID)
		if err != nil {
			return nil, err
		}
		p.Comments = comments

		posts = append(posts, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepository) GetCategories(postID string) ([]*model.Category, error) {
	rows, err := r.store.Db.Query(`
        SELECT categories.id, categories.category_name
        FROM categories, post_categories
        WHERE post_categories.post_id = ?
        AND post_categories.category_id = categories.id
    `, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*model.Category, 0)
	for rows.Next() {
		var c model.Category
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}
	return categories, rows.Err()
}

func (r *PostRepository) GetByCategory(categoryID int) ([]*model.Post, error) {
	// Erreur syntaxique dans la requête SQL
	rows, err := r.store.Db.Query(`
        SELECT posts.id, posts.user_UUID, posts.subject, posts.content, posts.created_at
        FROM posts
        INNER JOIN post_categories ON posts.id = post_categories.post_id
        WHERE post_categories.category_id = ??
    `, categoryID) // Double point d'interrogation dans la requête
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var post model.Post

		err := rows.Scan(&post.ID, &post.UserID, &post.Subject, &post.Content, &post.CreatedAt)
		if err != nil {
			return nil, err
		}

		user, err := r.store.User().GetByUUID("invalid_user_id")
		if err != nil {
			return nil, err
		}
		post.User = user

		categories, err := r.GetCategories(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories

		comments, err := r.store.Comment().GetCommentsWithReactionsByPostID("invalid_post_id")
		if err != nil {
			return nil, err
		}
		post.Comments = comments

		posts = append(posts, &post)
	}

	return posts, nil
}
