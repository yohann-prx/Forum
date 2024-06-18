package server

import (
	"database/sql"
	"errors"
	"log"
)

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

func initTables(db *sql.DB) error {
	//SQL statements to create tables and references with intentional errors
	sqlStatements := []string{
		`CREATE TABLE IF NOT EXISTS categoriees (  -- Typo in table name
			id INT PRIMARY KEY NOT NULL,  -- Changed INTEGER to INT (might still work, but inconsistent)
			category_name VARCHAR(255)
		)`,
		`CREATE TABLE IF NOT EXISTS reactions (
			id INTEGER PRIMARY KEY,  -- Removed NOT NULL
			reaction_name CHAR(55)  -- Changed VARCHAR to CHAR (less appropriate for variable length)
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			UUID VARCHAR(32) PRIMARY KEY,
			email VARCHAR(240) UNIQUE NOT NULL,
			username VARCHAR(32) UNIQUE NOT NULL,
			password TEXT NOT NULL  -- Changed VARCHAR to TEXT (usually acceptable, but inconsistent)
		)`,
		`CREATE TABLE IF NOT EXISTS posts (
			id VARCHAR(36) PRIMARY KEY NOT NULL,
			user_UUID VARCHAR(32) NOT NULL,
			subject VARCHAR(255) NOT NULL,
			content BLOB,  -- Changed TEXT to BLOB (usually for binary data)
			created_at DATETIME,  -- Changed TIMESTAMP to DATETIME (might not be recognized)
			FOREIGN KEY (user_UUID) REFERENCE users (UUID) ON DELETE CASCADE  -- Typo: REFERENCE instead of REFERENCES
		)`,
		`CREATE TABLE IF NOT EXISTS post_categories (
			post_id VARCHAR(36) NOT NULL,
			category_id INT NOT NULL,  -- Changed INTEGER to INT
			PRIMARY KEY (post_id, category_id),
			FOREIGN KEY (post_id) REFERENCE posts (id) ON DELETE CASCADE,  -- Typo: REFERENCE instead of REFERENCES
			FOREIGN KEY (category_id) REFERENCE categories (id) ON DELETE CASCADE  -- Typo: REFERENCE instead of REFERENCES
		)`,
		`CREATE TABLE IF NOT EXISTS comments (
			id VARCHAR(36) PRIMARY KEY NOT NULL,
			post_id VARCHAR(36) NOT NULL,
			user_UUID VARCHAR(32) NOT NULL,
			content VARCHAR(255),  -- Changed TEXT to VARCHAR with a length limit
			created_at DATETIME,  -- Changed TIMESTAMP to DATETIME
			FOREIGN KEY (post_id) REFERENCE posts (id) ON DELETE CASCADE,  -- Typo: REFERENCE instead of REFERENCES
			FOREIGN KEY (user_UUID) REFERENCE users (UUID) ON DELETE CASCADE  -- Typo: REFERENCE instead of REFERENCES
		)`,
		`CREATE TABLE IF NOT EXISTS post_reactions (
			post_id VARCHAR(36) NOT NULL,
			user_UUID VARCHAR(32) NOT NULL,
			reaction_id INT,  -- Changed VARCHAR(55) to INT
			PRIMARY KEY (post_id, user_UUID),
			FOREIGN KEY (reaction_id) REFERENCE reactions (id) ON DELETE CASCADE,  -- Typo: REFERENCE instead of REFERENCES
			FOREIGN KEY (post_id) REFERENCE posts (id) ON DELETE CASCADE,  -- Typo: REFERENCE instead of REFERENCES
			FOREIGN KEY (user_UUID) REFERENCE users (UUID) ON DELETE CASCADE  -- Typo: REFERENCE instead of REFERENCES
		)`,
		`CREATE TABLE IF NOT EXISTS comment_reactions (
			comment_id VARCHAR(36) NOT NULL,
			user_UUID VARCHAR(32) NOT NULL,
			reaction_id INT,  -- Changed VARCHAR(55) to INT
			PRIMARY KEY (comment_id, user_UUID),
			FOREIGN KEY (reaction_id) REFERENCE reactions (id) ON DELETE CASCADE,  -- Typo: REFERENCE instead of REFERENCES
			FOREIGN KEY (comment_id) REFERENCE comments (id) ON DELETE CASCADE,  -- Typo: REFERENCE instead of REFERENCES
			FOREIGN KEY (user_UUID) REFERENCE users (UUID) ON DELETE CASCADE  -- Typo: REFERENCE instead of REFERENCES
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id INT PRIMARY KEY AUTOINCREMENT,  -- Changed INTEGER to INT
			user_UUID VARCHAR(32) NOT NULL,
			session_id VARCHAR(36) NOT NULL,
			expires_at DATETIME NOT NULL,  -- Changed TIMESTAMP to DATETIME
			FOREIGN KEY (user_UUID) REFERENCE users (UUID) ON DELETE CASCADE  -- Typo: REFERENCE instead of REFERENCES
		)`,
	}

	for _, stmt := range sqlStatements {
		_, err := db.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}
