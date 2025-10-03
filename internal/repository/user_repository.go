package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func NewUserRepository(db *sql.DB) user.UserRepository {
	return &SQLiteUserRepository{DB: db}
}

// Create implements user.UserRepository.
func (rep *SQLiteUserRepository) Create(ctx context.Context, u *user.User, credentials *user.Credentials) (*user.User, error) {
	query := "INSERT INTO users (user_name, user_email, user_password) VALUES (?, ?, ?)"

	res, err := rep.DB.ExecContext(ctx, query, u.Name, u.Email, credentials.PasswordHash)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code() == 2067 {
			return nil, ErrDuplicateEmail
		}

		return nil, fmt.Errorf("não foi possível inserir o usuário no banco - %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("não foi possivel obter o ID do último usuário inserido - %w", err)
	}

	// Busca o user criado antes para retorná-lo na função
	u, err = rep.FindByID(ctx, id)
	if errors.Is(err, ErrUsrNotFound) {
		return nil, ErrUsrNotFound
	}
	if err != nil {
		return nil, err
	}

	return u, nil
}

// FindByEmail implements user.UserRepository.
func (rep *SQLiteUserRepository) FindByEmail(ctx context.Context, em string) (*user.User, *user.Credentials, error) {
	query := "SELECT user_id, user_name, user_email, user_password, user_created_at FROM users WHERE user_email = ?"
	var (
		userID        int64
		userName      string
		userEmail     string
		userPassHash  string
		userCreatedAt time.Time
	)

	err := rep.DB.QueryRowContext(ctx, query, em).Scan(&userID, &userName, &userEmail, &userPassHash, &userCreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, ErrUsrNotFound
	}
	if err != nil {
		return nil, nil, fmt.Errorf("não foi possível buscar o usuário - %w", err)
	}

	// Retornar o usuário e credenciais
	u := user.User{
		ID:        userID,
		Name:      userName,
		Email:     userEmail,
		CreatedAt: userCreatedAt,
	}
	password := user.Credentials{
		PasswordHash: userPassHash,
	}
	
	return &u, &password, nil
}

// FindByID implements user.UserRepository.
func (rep *SQLiteUserRepository) FindByID(ctx context.Context, id int64) (*user.User, error) {
	query := "SELECT user_id, user_name, user_email, user_created_at FROM users WHERE user_id = ?" 
	var (
		userID        int64
		userName      string
		userEmail     string
		userCreatedAt time.Time
	)

	err := rep.DB.QueryRowContext(ctx, query, id).Scan(&userID, &userName, &userEmail, &userCreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUsrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("não foi possível buscar o usuário - %w", err)
	}

	u := user.User {
		ID: userID,
		Name: userName,
		Email: userEmail,
		CreatedAt: userCreatedAt,
	}
	return &u, nil

}

