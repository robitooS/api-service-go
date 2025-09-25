package user

import (
	"fmt"
	"net/mail"
	"regexp"
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

func ValidatePassword(password string) error {
	// minimo de 1 letra maiscula, minuscula e 1 numero (minimo de 8 caracteres)
	// contem somente letras e numeros
	// bem xoxo esse regex aqui, da p fazer um melhor depois 
	var regex = regexp.MustCompile(`^(?:.*[a-z])(?:.*[A-Z])(?:.*\d)[a-zA-Z\d]{8,}$`)
	if !regex.MatchString(password) {
		return fmt.Errorf("senha inválida - a senha está fraca")
	}
	return nil
}