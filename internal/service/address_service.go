package service

import (
	"context"

	"github.com/robitooS/api-service-go/internal/domain/address"
)

type AddressService struct {
	AddressRepository address.AddressRepository
}

func NewAddressService(addressRepository  address.AddressRepository) *AddressService {
	return &AddressService{AddressRepository: addressRepository }
}

func (s *AddressService) Create(ctx context.Context, street, number, neighborhood, city, state, cep string, userID int64) (*address.Address, error) {
	add, err := address.NewAddress(street, number, neighborhood, city, state, cep, userID)
	if err != nil {
		return nil, err
	}

	addCreated, err := s.AddressRepository.Create(ctx, add)
	if err != nil {
		return nil, err
	}

	return addCreated, nil
}

func (s *AddressService) Update(ctx context.Context, userID int64, newAddress *address.Address) (*address.Address, error) {
	addressUpdated, err := s.AddressRepository.Update(ctx, userID, newAddress)
	if err != nil {
		return nil, err
	}

	return addressUpdated, nil
}