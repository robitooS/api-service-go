package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/robitooS/api-service-go/internal/auth"
	"github.com/robitooS/api-service-go/internal/cache"
	"github.com/robitooS/api-service-go/internal/domain/user"
	"github.com/robitooS/api-service-go/internal/handlers"
)

func AddressRoutes(router *gin.Engine, userRepository user.UserRepository, addressHandler *handlers.AddressHandler, cache cache.NonceStore, hmacKey []byte) {

	privateGroup := router.Group("/address", auth.AuthenticateHMAC(hmacKey, userRepository, cache))
	{
		privateGroup.POST("/create", addressHandler.CreateAddress)
		privateGroup.POST("/update", addressHandler.UpdateAddress)
	}
}