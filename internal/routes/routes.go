package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/robitooS/api-service-go/internal/cache"
	"github.com/robitooS/api-service-go/internal/domain/user"
	"github.com/robitooS/api-service-go/internal/handlers"
)

func SetupRoutes( r *gin.Engine, userRepository user.UserRepository, userHandler *handlers.UserHandler, addressHandler *handlers.AddressHandler, cache cache.NonceStore,  hmacKey []byte) {
	UserRoutes(r, userRepository, userHandler, cache, hmacKey)
	AddressRoutes(r, userRepository, addressHandler, cache, hmacKey)
}