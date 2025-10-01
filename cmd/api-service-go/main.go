package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/robitooS/api-service-go/internal/cache"
	"github.com/robitooS/api-service-go/internal/config"
	"github.com/robitooS/api-service-go/internal/database"
	"github.com/robitooS/api-service-go/internal/routes"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("erro ao carregar config: %v", err)
	}

	pool, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("erro ao conectar no banco: %v", err)
	}
	defer pool.Close()

	if err := db.RunMigrations(pool); err != nil {
		log.Fatalf("erro ao rodar migrations: %v", err)
	}
	fmt.Println("Migrations executadas com sucesso.")

	router := gin.Default()
	cache := cache.NewInMemoryNonceStore()
	routes.SetupRoutes(router, pool, cfg, cache)

	fmt.Printf("[INFO] Servidor configurado e escutando na porta %s\n", cfg.HttpAddr)
	router.Run(cfg.HttpAddr)
}