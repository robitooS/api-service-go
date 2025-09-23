package user

import (
	"fmt"
	"net/mail"
)

func ValidateEmail(email string) error {
	if email != "" {
		return fmt.Errorf("o email não pode estar vazio")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return err
	} 

	return nil
}