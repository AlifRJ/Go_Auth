package migration

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func UserMigration(ctx context.Context, db *pgxpool.Pool) error {
	query := `
	-- 1. Tabel Users
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		username VARCHAR(255) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at TIMESTAMP,
		deleted_at TIMESTAMP
	);

	-- 2. Index untuk pencarian cepat Email/Username pada Login & Soft Delete
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE deleted_at IS NULL;
	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username) WHERE deleted_at IS NULL;

	-- 3. Tabel Auth Sessions (Refresh Token Persistence)
	CREATE TABLE IF NOT EXISTS auth_sessions (
		id SERIAL PRIMARY KEY,
		user_id INT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		refresh_token TEXT NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
	);

	-- 4. Index untuk pencarian sesi berdasarkan user_id
	CREATE INDEX IF NOT EXISTS idx_auth_sessions_user_id ON auth_sessions(user_id);
	`

	if _, err := db.Exec(ctx, query); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}