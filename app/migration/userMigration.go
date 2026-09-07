package migration

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func UserMigration(ctx context.Context, db *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		username VARCHAR(255) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NUll,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at TIMESTAMP,
		deleted_at TIMESTAMP
	);`

	if _, err := db.Exec(ctx, query); err != nil {
		return err
	}
	return nil
}