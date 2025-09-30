package service

import (
	"context"
	"fmt"
	"time"

	"github.com/robitooS/api-service-go/internal/auth"
	"github.com/robitooS/api-service-go/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepository user.UserRepository
	HmacSecret []byte
}

func NewUserService(userRepository user.UserRepository, hmacSecret string) *UserService {
	return &UserService{
		UserRepository: userRepository,
		HmacSecret: []byte(hmacSecret),
	}
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

// Struct de resposta após login
type AuthResponse struct {
	UserID int64 `json:"user_id"`
	Timestamp int64 `json:"timestamp"`
	Signature string `json:"signature"`
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

	timestamp := time.Now().Unix()
	msg := auth.BuildMessage(user.ID, timestamp)
	signature := auth.GenerateSignature(msg, s.HmacSecret)

	resp := AuthResponse{
		UserID: user.ID,
		Timestamp: timestamp,
		Signature: signature,
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