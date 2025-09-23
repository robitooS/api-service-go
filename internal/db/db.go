package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	_ "modernc.org/sqlite"
	"github.com/robitooS/api-service-go/internal/config"
)

//! A dsn será passado pela main
func Connect(cfg *config.Config) (*sql.DB, error) {
	var err error

	pool, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		// Aqui não é um erro de conexão, pode ser erro sobre o dsn ou outros 
		return nil, fmt.Errorf("não foi possível utilizar o dsn: %v", err)
	}

	// Seta limite de conexões do pool, tempo de vida da conn e qtd de conn ociosas guardadas 
	pool.SetConnMaxLifetime(0)
	pool.SetMaxIdleConns(3)
	pool.SetMaxOpenConns(3)

	// Validar a conexão do db com ping
	err = ping(context.Background(), pool)
	if err != nil {
		pool.Close() // Fechar o pool de conexões caso o ping falhe
		return nil, err
	}

	return pool, nil
}

func ping(ctx context.Context, pool *sql.DB) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := pool.PingContext(ctx); err != nil {
		return fmt.Errorf("não foi possível efetuar uma conexão com o db: %v", err)
	}
	return nil
}