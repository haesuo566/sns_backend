package entities

import "time"

type User struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UserTag   string    `json:"user_tag"`
	// UserKey   string    `json:"user_key"`
	Platform string `json:"platform"`
}
