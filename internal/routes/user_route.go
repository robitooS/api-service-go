package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/robitooS/api-service-go/internal/auth"
	"github.com/robitooS/api-service-go/internal/cache"
	"github.com/robitooS/api-service-go/internal/domain/user"
	"github.com/robitooS/api-service-go/internal/handlers"
)

func UserRoutes(router *gin.Engine, userRepository user.UserRepository, userHandler *handlers.UserHandler, cache cache.NonceStore, hmacKey []byte) {

	// Rotas públicas
	publicUsersRoutes := router.Group("/users")
	{
		publicUsersRoutes.POST("/create", userHandler.CreateUser)
		publicUsersRoutes.POST("/login", userHandler.Login)
	}

	// Rotas protegidas (HMAC)
	authUsersRoutes := router.Group("/users", auth.AuthenticateHMAC(hmacKey, userRepository, cache))
	{
		authUsersRoutes.POST("/get", userHandler.GetUserByID)
	}
}