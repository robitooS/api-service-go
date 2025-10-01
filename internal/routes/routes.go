package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/robitooS/api-service-go/internal/cache"
	"github.com/robitooS/api-service-go/internal/config"
)

func SetupRoutes(r *gin.Engine, pool *sql.DB, cfg *config.Config, cache cache.NonceStore) {
	// Router para users
	UserRoutes(r, pool, cfg, cache)
}