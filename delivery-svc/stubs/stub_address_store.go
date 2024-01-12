package stubs

import (
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/storeerrors"
)

type StubAddressStore struct {
	CreatedAddresses []models.Address
	Addresses        []models.Address
}

func (s *StubAddressStore) CreateAddress(address *models.Address) error {
	s.CreatedAddresses = append(s.CreatedAddresses, *address)
	address.ID = len(s.CreatedAddresses)

	return nil
}

func (s *StubAddressStore) GetAddressByID(id int) (models.Address, error) {
	for _, address := range s.Addresses {
		if address.ID == id {
			return address, nil
		}
	}
	return models.Address{}, storeerrors.ErrNotFound
}
