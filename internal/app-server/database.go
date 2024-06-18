package server

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the database with tables and returns the *sql.DB object.
func InitDB(config Config) (*sql.DB, error) {
	// Open database connection
	db, err := sql.Open(config.DbDriver, config.DbPath)
	if err != nil {
		return nil, errors.Join(errors.New("error initializing database"), err)
	}

	// Initialize tables
	if err = initTables(db); err != nil {
		return nil, err
	}

	// Check database connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Log success message
	log.Println("SQLite database initialized successfully.")
	// Return *sql.DB object
	return db, nil
}

// initTables creates necessary tables in the database.
func initTables(db *sql.DB) error {
	//SQL statements to create tables and references
	sqlStatements := []string{
		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY NOT NULL,
			category_name VARCHAR(255)
		)`,
		`CREATE TABLE IF NOT EXISTS reactions (
			id INTEGER PRIMARY KEY NOT NULL,
			reaction_name VARCHAR(55)
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			UUID VARCHAR(32) PRIMARY KEY,
			email VARCHAR(240) UNIQUE NOT NULL,
			username VARCHAR(32) UNIQUE NOT NULL,
			password VARCHAR(240) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS posts (
			id VARCHAR(36) PRIMARY KEY NOT NULL,
			user_UUID VARCHAR(32) NOT NULL,
			subject VARCHAR(255) NOT NULL,
			content TEXT,
			created_at TIMESTAMP,
			FOREIGN KEY (user_UUID) REFERENCES users (UUID) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS post_categories (
			post_id VARCHAR(36) NOT NULL,
			category_id INTEGER NOT NULL,
			PRIMARY KEY (post_id, category_id),
			FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
			FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS comments (
			id VARCHAR(36) PRIMARY KEY NOT NULL,
			post_id VARCHAR(36) NOT NULL,
			user_UUID VARCHAR(32) NOT NULL,
			content TEXT,
			created_at TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
			FOREIGN KEY (user_UUID) REFERENCES users (UUID) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS post_reactions (
			post_id VARCHAR(36) NOT NULL,
			user_UUID VARCHAR(32) NOT NULL,
			reaction_id VARCHAR(55),
			PRIMARY KEY (post_id, user_UUID),
			FOREIGN KEY (reaction_id) REFERENCES reactions (id) ON DELETE CASCADE,
			FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
			FOREIGN KEY (user_UUID) REFERENCES users (UUID) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS comment_reactions (
			comment_id VARCHAR(36) NOT NULL,
			user_UUID VARCHAR(32) NOT NULL,
			reaction_id VARCHAR(55),
			PRIMARY KEY (comment_id, user_UUID),
			FOREIGN KEY (reaction_id) REFERENCES reactions (id) ON DELETE CASCADE,
			FOREIGN KEY (comment_id) REFERENCES comments (id) ON DELETE CASCADE,
			FOREIGN KEY (user_UUID) REFERENCES users (UUID) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_UUID VARCHAR(32) NOT NULL,
			session_id VARCHAR(36) NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			FOREIGN KEY (user_UUID) REFERENCES users (UUID) ON DELETE CASCADE
		)`,
	}

	// Execute SQL statements
	for _, statement := range sqlStatements {
		_, err := db.Exec(statement)
		if err != nil {
			return errors.Join(errors.New("error executing table initialization"), err)
		}
	}
	return nil
}
