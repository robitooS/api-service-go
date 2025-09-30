package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/robitooS/api-service-go/internal/config"
	"github.com/robitooS/api-service-go/internal/handlers"
	"github.com/robitooS/api-service-go/internal/repository"
	"github.com/robitooS/api-service-go/internal/service"
)

func UserRoutes(router *gin.Engine, pool *sql.DB, cfg *config.Config) {
	// primeiro deve injetar as dependências 
	userRepository := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService) // handler responsável pelo usuário

	// Rotas públicas
	publicUsersRoutes := router.Group("/users")
	{
		publicUsersRoutes.POST("list") // listar todos os usuários
		publicUsersRoutes.POST("/create", userHandler.CreateUser)
		publicUsersRoutes.POST("/get") // buscar por id especifico
		
		// daq p frente será adicionado mais rotas
	}

	// // Rotas protegidas (HMAC)
	// authUsersRoutes := router.Group("/users", auth.AuthenticateHMAC([]byte(cfg.HmacSecret), userRepository))
	// {

	// }
}