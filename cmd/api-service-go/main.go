package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/robitooS/api-service-go/internal/cache"
	"github.com/robitooS/api-service-go/internal/config"
	"github.com/robitooS/api-service-go/internal/database"
	"github.com/robitooS/api-service-go/internal/handlers"
	"github.com/robitooS/api-service-go/internal/repository"
	"github.com/robitooS/api-service-go/internal/routes"
	"github.com/robitooS/api-service-go/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("erro ao carregar config: %v", err)
	}
	hmacKey, err := base64.StdEncoding.DecodeString(cfg.HmacSecret)
	if err != nil {
		log.Fatalf("HMAC_SECRET inválido") // indica que não ta codificado com base64
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

	cache := cache.NewInMemoryNonceStore()

	// Repositórios
	userRepository := repository.NewUserRepository(pool)
	addressRepository := repository.NewAddressRepository(pool)

	// Serviços
	userService := service.NewUserService(userRepository)
	addressService := service.NewAddressService(addressRepository)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	addressHandler := handlers.NewAddressHandler(addressService)
	

	router := gin.Default()
	routes.SetupRoutes(router, userRepository, userHandler, addressHandler, cache, hmacKey)


	fmt.Printf("[INFO] Servidor configurado e escutando na porta %s\n", cfg.HttpAddr)
	router.Run(cfg.HttpAddr)
}