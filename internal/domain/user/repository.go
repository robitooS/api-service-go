package user

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User, credentials *Credentials) (*User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, *Credentials, error)
}
