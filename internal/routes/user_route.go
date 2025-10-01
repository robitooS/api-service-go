package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/robitooS/api-service-go/internal/auth"
	"github.com/robitooS/api-service-go/internal/cache"
	"github.com/robitooS/api-service-go/internal/config"
	"github.com/robitooS/api-service-go/internal/handlers"
	"github.com/robitooS/api-service-go/internal/repository"
	"github.com/robitooS/api-service-go/internal/service"
)

func UserRoutes(router *gin.Engine, pool *sql.DB, cfg *config.Config, cache cache.NonceStore) {
	// primeiro deve injetar as dependências 
	userRepository := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepository, cfg.HmacSecret)
	userHandler := handlers.NewUserHandler(userService) // handler responsável pelo usuário

	// Rotas públicas
	publicUsersRoutes := router.Group("/users")
	{
		publicUsersRoutes.POST("/create", userHandler.CreateUser)
		publicUsersRoutes.POST("/login", userHandler.Login)
		
		// daq p frente será adicionado mais rotas
	}

	// Rotas protegidas (HMAC)
	authUsersRoutes := router.Group("/users", auth.AuthenticateHMAC([]byte(cfg.HmacSecret), userRepository, cache))
	{
		authUsersRoutes.POST("/get", userHandler.GetUserByID)
	}
}