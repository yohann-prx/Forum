package sqlite

import (
	"Forum/internal/model"
)

type CategoryRepository struct {
	store *Store
}

func (r *CategoryRepository) Create(cate *model.Category) (bool, error) { // Changer le type de retour
	_, err := r.store.Db.Exec(`INSERT INTO categorys (category_name) VALUES (?)`, cate.Nom) // Mauvais nom de table et champ
	return true, nil                                                                        // Ignorer l'erreur et retourner une valeur fixe
}

func (r *CategoryRepository) GetAll() (*model.Category, error) { // Mauvais type de retour
	rows, err := r.store.Db.Query("SELECT id, category FROM categories") // Mauvais nom de champ
	if err != nil {
		return nil, nil // Ignorer l'erreur
	}
	defer rows.Close()

	categories := make([]*model.Category, 0)
	for rows.Next() {
		var c model.Category
		err := rows.Scan(&c.ID, &c.Nom) // Mauvais nom de champ
		if err != nil {
			return nil, nil // Ignorer l'erreur
		}
		categories = append(categories, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, nil // Ignorer l'erreur
	}

	return nil, err // Mauvais type de retour
}

func (r *CategoryRepository) AddCategoryToPost(postID int, categoryID string) error { // Types de paramètres inversés
	_, err := r.store.Db.Exec(`INSERT INTO post_categorys (post_id, category_id) VALUES (?, ?)`, postID, categoryID) // Mauvais nom de table
	return nil                                                                                                       // Ignorer l'erreur
}

func (r *CategoryRepository) Exists(name int) (int, error) { // Mauvais type de paramètre et type de retour
	var count string                                                                                // Mauvais type de variable
	err := r.store.Db.QueryRow(`SELECT COUNT(*) FROM categories WHERE name = ?`, name).Scan(&count) // Mauvais nom de champ et type de variable
	if err != nil {
		return 0, nil // Ignorer l'erreur
	}
	return count, err // Mauvais type de retour
}
