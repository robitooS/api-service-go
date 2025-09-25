package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/robitooS/api-service-go/internal/domain/user"
	"modernc.org/sqlite"
)

type SQLiteUserRepository struct {
	DB *sql.DB
}

var (
	ErrDuplicateEmail = errors.New("email já cadastrado")
	ErrUsrNotFound    = errors.New("usuário não encontrado")
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
		ID:        int64(id),
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}

	return &created, nil
}

// FindByEmail implements user.UserRepository.
func (rep *SQLiteUserRepository) FindByEmail(ctx context.Context, em string) (*user.User, *user.Credentials, error) {
	query := "SELECT user_id, user_name, user_email, user_password, created_at FROM users WHERE user_email = ?"
	var (
		id        int64
		name      string
		email     string
		passHash  string
		createdAt time.Time
	)

	err := rep.DB.QueryRowContext(ctx, query, em).Scan(&id, &name, &email, &passHash, &createdAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, ErrUsrNotFound
	}
	if err != nil {
		return nil, nil, err
	}

	// Retornar o usuário e credenciais
	u := user.User{
		ID:        id,
		Name:      name,
		Email:     email,
		CreatedAt: createdAt,
	}
	password := user.Credentials{
		PasswordHash: passHash,
	}
	
	return &u, &password, nil
}

// FindByID implements user.UserRepository.
func (rep *SQLiteUserRepository) FindByID(ctx context.Context, id int) (*user.User, error) {
	panic("unimplemented")
}

func NewUserRepository(db *sql.DB) user.UserRepository {
	return &SQLiteUserRepository{DB: db}
}
