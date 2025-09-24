package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID    int64
	Name  string
	Email string
	CreatedAt time.Time
}

func NewUser(name, email, password string) (*User, *Credentials, error) {
	if err := ValidateName(name); err != nil {
		return nil, nil, err
	}
	if err := ValidateEmail(email); err != nil {
		return nil, nil, err
	}
	if err := ValidatePassword(password); err != nil {
		return nil, nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	return &User{Name: name, Email: email}, &Credentials{PasswordHash: string(hash)}, nil
}
