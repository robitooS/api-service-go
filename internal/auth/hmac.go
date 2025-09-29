package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

// Vai gerar o novo HMAC conforme a mensagem recebida
func GenerateSignature(msg string, key []byte) string {
	// Criar código hash
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(msg))
	
	// Decodifica slice de bytes p retornar uma string legivel
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return signature
}

// Valida se a assinatura do hmac vindo bate com a key do servidor
func ValidateSignature(msg, signature string, key []byte) (bool, error) {
	expectedHMAC := hmac.New(sha256.New, key)
	expectedHMAC.Write([]byte(msg))
	// Chave esperada de acordo com a mensagem
	expectedSign := expectedHMAC.Sum(nil)

	// Assinatura vinda da mensagem
	sign, err := base64.RawURLEncoding.DecodeString(signature); 
	
	if err != nil {
		return false, fmt.Errorf("não foi possível decodificar a assinatura vinda da requisição - %w", err)
	}
	
	return hmac.Equal(sign, expectedSign), nil
}

func ValidateTimeStamp(ts int64) error {
	window := int64(300)
	now := time.Now().Unix()

	if ts > now + window { // timestamp no futuro
		return fmt.Errorf("o timestamp passou do limite esperado")
	}

	if ts < now - window { // timestamp muito antigo
		return fmt.Errorf("o timestamp está muito atrasado")
	}

	return nil
}

func BuildMessage(userID, ts int64) string {
	return fmt.Sprintf("%d:%d", userID, ts)
}