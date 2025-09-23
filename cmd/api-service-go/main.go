package main

import (
	"fmt"
	"log"

	"github.com/robitooS/api-service-go/internal/config"
	"github.com/robitooS/api-service-go/internal/db"
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
}