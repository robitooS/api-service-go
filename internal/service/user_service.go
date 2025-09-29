package service

import (
	"context"

	"github.com/robitooS/api-service-go/internal/domain/user"
)

type UserService struct {
	UserRepository user.UserRepository
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
