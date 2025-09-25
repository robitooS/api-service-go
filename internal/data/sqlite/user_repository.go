package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/robitooS/api-service-go/internal/domain/user"
	"modernc.org/sqlite"
)

type SQLiteUserRepository struct {
	DB *sql.DB
}

var (
	ErrDuplicateEmail = errors.New("email já cadastrado")
)

// Create implements user.UserRepository.
func (rep *SQLiteUserRepository) Create(ctx context.Context, u *user.User, credentials *user.Credentials) (*user.User, error) {
	query := "INSERT INTO users (user_name, user_email, user_password) VALUES (?, ?, ?)"
	
	res, err := rep.DB.ExecContext(ctx, query, u.Name, u.Email, credentials.PasswordHash)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code() == 2067 {
			return nil, ErrDuplicateEmail
		}

		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	created := user.User{
		ID: int64(id),
		Name: u.Name,
		Email: u.Email,
	}

	return &created, nil
}

// FindByEmail implements user.UserRepository.
func (rep *SQLiteUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, *user.Credentials, error) {
	panic("unimplemented")
}

// FindByID implements user.UserRepository.
func (rep *SQLiteUserRepository) FindByID(ctx context.Context, id int) (*user.User, error) {
	panic("unimplemented")
}

func NewUserRepository(db *sql.DB) user.UserRepository {
	return &SQLiteUserRepository{DB: db}
}
