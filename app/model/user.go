package model

import (
	"context"
	"time"
)

type User struct {
	ID         uint			`json:"id"`
	Name       string		`json:"name"`
	Username   string		`json:"username"`
	Email      string		`json:"email"`
	Password   string		`json:"-"`
	Created_at time.Time	`json:"created_at"`
	Updated_at *time.Time	`json:"updated_at"`
	Deleted_at *time.Time	`json:"deleted_at"`
}

type UserRepository interface {
	Login(ctx context.Context, identity string) (*User, error)
	GetAll(ctx context.Context, limit, offset int) ([]*User, error)
	GetByID(ctx context.Context, id uint) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
	PermanentDelete(ctx context.Context, id uint) error
}