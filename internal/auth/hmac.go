package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Vai gerar o novo HMAC conforme a mensagem recebida
func GenerateSignature(msg string, key []byte) (string, error) {
	// Criar código hash
	mac := hmac.New(sha256.New, key)

	if _, err := mac.Write([]byte(msg)); err != nil {
		return "", fmt.Errorf("não foi possível adicionar o conteúdo da mensagem a assinatura - %w", err)
	}

	// Decodifica slice de bytes p retornar uma string legivel
	signature := hex.EncodeToString(mac.Sum(nil))
	return signature, nil
}

// Valida se a assinatura do hmac vindo bate com a key do servidor
func ValidateSignature(msg, signature, key string) bool {
	return false
}