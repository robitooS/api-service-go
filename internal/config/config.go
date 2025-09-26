package config

import (
	"errors"
	"fmt"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv     string
	HttpAddr   string
	DBPath     string
	HmacSecret string
}

func Load() (*Config, error){ 
	var configs Config 
	err := godotenv.Load()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			
			// ******************************************
			// Arquivo .env não encontrado, modo de prod.
			// continuar fluxo
			// ******************************************

		} else {
			// Erro de permissão, corrompido, etc etc ...
			return nil, fmt.Errorf("não foi possível carregar o .env: %v", err)
		}
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "dev"
	}

	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "internal/data/app.db"
	}

	hmacSecret := os.Getenv("HMAC_SECRET")
	if hmacSecret == "" {
		return nil, fmt.Errorf("key hmac_secret não pode estar vazia")
	}

	configs.AppEnv = appEnv
	configs.HttpAddr = httpAddr
	configs.DBPath = dbPath
	configs.HmacSecret = hmacSecret

	return &configs, nil
}
