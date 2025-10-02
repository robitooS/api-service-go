package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/robitooS/api-service-go/internal/cache"
)

func SetupRoutes(r *gin.Engine, pool *sql.DB, cache cache.NonceStore, hmacKey []byte) {
	// Router para users
	UserRoutes(r, pool, cache, hmacKey)
}