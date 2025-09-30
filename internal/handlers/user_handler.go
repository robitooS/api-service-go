package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/robitooS/api-service-go/internal/service"
)

type UserHandler struct {
	UserService *service.UserService
}

type CreateUserRequest struct {
	Name string `json:"user_name"`
	Email string `json:"user_email"`
	Password string `json:"user_password"`
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	request := CreateUserRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error":fmt.Sprintf("body da requisição inválido: %s", err)})
		return
	}

	user, err := uh.UserService.Create(ctx.Request.Context(), request.Name, request.Email, request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":fmt.Sprintf("falha ao persistir usuário no banco: %s", err)})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}