package testutil

import (
	"fmt"

	"github.com/VitoNaychev/bt-customer-svc/models"
	td "github.com/VitoNaychev/bt-customer-svc/testdata"
)

type StubAddressStore struct {
	addresses   []models.Address
	storeCalls  []models.Address
	deleteCalls []int
	updateCalls []models.Address
}

func NewStubAddressStore(data []models.Address) *StubAddressStore {
	return &StubAddressStore{
		addresses:   data,
		storeCalls:  []models.Address{},
		deleteCalls: []int{},
		updateCalls: []models.Address{},
	}
}

func (s *StubAddressStore) GetAddressesByCustomerId(customerId int) ([]models.Address, error) {
	if customerId == td.PeterCustomer.Id {
		return []models.Address{td.PeterAddress1, td.PeterAddress2}, nil
	}

	if customerId == td.AliceCustomer.Id {
		return []models.Address{td.AliceAddress}, nil
	}

	return []models.Address{}, nil
}

func (s *StubAddressStore) StoreAddress(address models.Address) {
	s.storeCalls = append(s.storeCalls, address)
}

func (s *StubAddressStore) DeleteAddressById(id int) error {
	_, err := s.GetAddressById(id)
	if err != nil {
		return err
	} else {
		s.deleteCalls = append(s.deleteCalls, id)
		return nil
	}
}

func (s *StubAddressStore) GetAddressById(id int) (models.Address, error) {
	for _, address := range s.addresses {
		if address.Id == id {
			return address, nil
		}
	}
	return models.Address{}, fmt.Errorf("address with id %d doesn't exist", id)
}

func (s *StubAddressStore) UpdateAddress(address models.Address) error {
	s.updateCalls = append(s.updateCalls, address)
	return nil
}

func (s *StubAddressStore) Empty() {
	s.addresses = []models.Address{}
	s.storeCalls = []models.Address{}
}
