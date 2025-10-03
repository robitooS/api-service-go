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

func NewAddress(id int64, address string, number string, neighborhood string, city string, state string, cep string, userID int64) (*Address, error) {
	
	return &Address{
		ID: id,
		Address: address,
		Number: number,
		Neighborhood: neighborhood,
		City: city,
		State: state,
		CEP: cep,
		UserID: userID,
	}, nil
}