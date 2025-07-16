package models

import "time"

type User struct {
	ID           int64     `json:"id"`
	Login        string    `json:"login"`
	Password     string    `json:"-"` // password should not be sent in JSON
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}
