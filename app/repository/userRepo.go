package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlifRJ/Go_Auth/app/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) model.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Login(ctx context.Context, identity string) (*model.User, error) {
	query := `SELECT id, name, username, email, password FROM users WHERE username = $1 OR email = $1`
	
	var u model.User
	err := r.db.QueryRow(ctx, query, identity).Scan(&u.ID, &u.Name, &u.Username, &u.Email, &u.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	
	return &u, nil
}

func (r *PostgresUserRepository) GetAll(ctx context.Context, limit, offset int) ([]*model.User, error) {
	query := `SELECT id, name, username, email, created_at, updated_at, deleted_at FROM users ORDER BY id DESC LIMIT $1 OFFSET $2`
	
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		fmt.Println(err)
		return nil, err
	} 
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Username, &u.Email, &u.Created_at, &u.Updated_at, &u.Deleted_at); err != nil {
			fmt.Println(err)
			return nil, err
		}
		users = append(users, &u)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Return [] if nil
	if users == nil {
		users = make([]*model.User, 0)
	}
	
	return users, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	query := `SELECT id, name, username, email, created_at, updated_at, deleted_at FROM users WHERE id = $1`
	
	var u model.User
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Name, &u.Username, &u.Email, &u.Created_at, &u.Updated_at, &u.Deleted_at)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	
	return &u, nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (name, username, email, password) VALUES ($1, $2, $3, $4) RETURNING id`
	
	err := r.db.QueryRow(ctx, query, user.Name, user.Username, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		return err
	}
	
	return nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *model.User) error {
	query := `UPDATE users SET name=$1, username=$2, email=$3, password=$4, updated_at=CURRENT_TIMESTAMP WHERE id = $5`
	
	tag, err := r.db.Exec(ctx, query, user.Name, user.Username, user.Email, user.Password, user.ID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0{
		return errors.New("User not found to update")
	}
	
	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id uint) error {
	query := `UPDATE users SET deleted_at = CURRENT_TIMESTAMP where id = $1`
	
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0{
		return errors.New("User not found to delete")
	}
	
	return nil
}

func (r *PostgresUserRepository) PermanentDelete(ctx context.Context, id uint) error {
	query := `DELETE FROM users WHERE id = $1`
	
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	
	if commandTag.RowsAffected() == 0 {
		return errors.New("User not found to delete")
	}
	
	return nil
}
