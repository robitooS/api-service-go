package address

import "time"

type Address struct {
	ID int64
	Address string
	Number string
	Neighborhood string
	City string
	State string
	CEP string
	UserID int64
	CreatedAt time.Time
}

func NewAddress(address string, number string, neighborhood string, city string, state string, cep string, userID int64) (*Address, error) {
	fields := []struct {
		name string
		value string
		max int
		min int
	} {
		{"endereço", address, 200, 2},
		{"numero", number, 10, 1},
		{"bairro", neighborhood, 80, 2},
		{"cidade", city, 100, 2},
		{"estado", state, 2, 2},
		{"cep", cep, 9, 9},
	}
	
	// Validações básiscas 
	for _, field := range fields {
		if err := ValidateRequiredField(field.name, field.value, field.max, field.min); err != nil {
			return nil, err
		}
	}
	
	// Validações de formatações
	if err := ValidateState(state); err != nil {
		return nil, err
	}

	if err := ValidateCEP(cep); err != nil {
		return nil, err
	}

	return &Address{
		Address: address,
		Number: number,
		Neighborhood: neighborhood,
		City: city,
		State: state,
		CEP: cep,
		UserID: userID,
	}, nil
}