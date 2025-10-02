package user

import (
	"fmt"
	"net/mail"
	"strings"
)

func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email inválido - o endereço está vazio")
	}
	if len(email) > 255 {
		return fmt.Errorf("email inválido - o tamanho não pode ultrapassar 255 caracteres")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("email inválido - formato incorreto")
	} 

	return nil
}

func ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("nome inválido - o nome está vazio")
	}
	if len(name) < 2 {
		return fmt.Errorf("nome inválido - o nome deve possuir mais de 1 caractere")
	}
	if len(name) > 100 {
		return fmt.Errorf("nome inválido - o nome não deve ultrapassar 100 caracteres")
	}

	return nil
}

// teve q trocar o regex q tava tudo errado, nãu funcionava o lookahead
// a verificação é feita manual agora
func ValidatePassword(password string) error {
	// senha precisa ter no mínimo 8 caracteres
	if len(password) < 8 {
		return fmt.Errorf("senha inválida - mínimo de 8 caracteres")
	}

	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, c := range password {
		switch {
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= '0' && c <= '9':
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:,.<>?/`~", c):
			hasSpecial = true
		default:
			// qualquer outro caractere não é permitido
			return fmt.Errorf("senha inválida - contém caractere não permitido")
		}
	}

	if !hasLower || !hasUpper || !hasDigit || !hasSpecial {
		return fmt.Errorf("senha inválida - precisa ter maiúscula, minúscula, número e símbolo especial")
	}
	return nil
}