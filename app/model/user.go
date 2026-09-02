package model

import "time"

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