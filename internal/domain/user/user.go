package user

import "time"

type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}