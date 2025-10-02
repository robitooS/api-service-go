package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type NonceStore interface {
	CacheNonce(nonce string) error
	runCleaner()
}

type InMemoryNonceStore struct {
	Storage map[string]time.Time
	mu		sync.RWMutex
}

// CacheNonce implements NonceStore.
func (i *InMemoryNonceStore) CacheNonce(nonce string) error {
	// evitar concorrencia de escrita
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.Storage[nonce]; ok {
		log.Println("[INFO-NONCE] O NONCE JÁ FOI USADO")
		return fmt.Errorf("o nonce já foi usado")

	}
	log.Println("[INFO-NONCE] NONCE DISPONÍVEL")
	i.Storage[nonce] = time.Now()
	return nil
}

func NewInMemoryNonceStore() NonceStore {
	// Cria o objeto
	inMemoryNonceStore := &InMemoryNonceStore{Storage: make(map[string]time.Time)}

	go inMemoryNonceStore.runCleaner()

	return inMemoryNonceStore
}

func (i *InMemoryNonceStore) runCleaner() {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        i.mu.Lock()

		// Ficar verificando se já expirou os 5 minutos
		// Vamo ficar vendo o valor de cada nonce e vendo se já passou 5 minutos
		// time desde o tempo colocado no map é maior q 5 minutos?
        for nonce, timestamp := range i.Storage {
            if time.Since(timestamp) > 5*time.Minute {
                delete(i.Storage, nonce)
            }
        }
        i.mu.Unlock()
    }
}