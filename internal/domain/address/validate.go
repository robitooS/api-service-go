package address

import (
	"fmt"
	"regexp"
)

func ValidateRequiredField(fieldName, value string, maxLength int, minLenght int) error {
	if value == "" {
		return fmt.Errorf("%s inválido(a) - o campo não pode estar vazio", fieldName)
	}
	if len(value) > maxLength {
		return fmt.Errorf("%s inválido(a) - o campo não deve ultrapassar %d caracteres", fieldName, maxLength)
	}
	if len(value) < minLenght {
		return fmt.Errorf("%s inválido(a) - o campo deve conter no mìnimo %d caracteres", fieldName, minLenght)
	}
	return nil
}

func ValidateState(state string) error {
	if err := ValidateRequiredField("estado", state, 2, 2); err != nil {
		return err
	}
	
	if ok, _ := regexp.MatchString(`^[a-zA-Z]+$`, state); !ok {
		return fmt.Errorf("estado inválido - deve conter apenas letras")
	}
	return nil
}

// tem q tar no formato 00000-000
func ValidateCEP(cep string) error {
	if err := ValidateRequiredField("cep", cep, 9, 2); err != nil {
		return err
	}
	
	re := regexp.MustCompile(`^\d{5}-\d{3}$`)
	if !re.MatchString(cep) {
		return fmt.Errorf("cep inválido - o formato esperado é XXXXX-XXX")
	}
	return nil
}