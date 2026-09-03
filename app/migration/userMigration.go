package migration

import (
	"database/sql"
)

func UserMigration(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		username VARCHAR(255) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NUll,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP,
		deleted_at TIMESTAMP
	);`

	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}