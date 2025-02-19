package domain

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // パスワードはJSON出力しない
	CreatedAt time.Time `json:"created_at"`
}
