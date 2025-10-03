package service

import (
	"context"
	"fmt"

	"github.com/robitooS/api-service-go/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepository user.UserRepository
}

// Struct de resposta após login
type AuthResponse struct {
	UserID int64 `json:"user_id"`
}

func NewUserService(userRepository user.UserRepository) *UserService {
	return &UserService{UserRepository: userRepository}
}

func (s *UserService) Create (ctx context.Context, name string, email string, password string) (*user.User, error){
	// Aqui vamo persistir o usuário no banco
	// Primeiro criar o usuário 
	user, credentials, err := user.NewUser(name, email, password)
	if err != nil {
		return nil, err
	}

	// User persistido no banco
	createdUser, err := s.UserRepository.Create(ctx, user, credentials)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *UserService) Login (ctx context.Context, email string, password string) (*AuthResponse, error) {
	// Verifica primeiro se o user existe
	user, creds, err := s.UserRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("usuário não encontrado")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(creds.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("credenciais inválidas")
	}

	resp := AuthResponse{
		UserID: user.ID,
	}

	return &resp, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*user.User, error) {
	user, err := s.UserRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}