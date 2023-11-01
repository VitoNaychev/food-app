package testutil

import (
	"fmt"

	as "github.com/VitoNaychev/bt-customer-svc/models/address_store"
	td "github.com/VitoNaychev/bt-customer-svc/testdata"
)

type StubAddressStore struct {
	addresses   []as.Address
	storeCalls  []as.Address
	deleteCalls []int
	updateCalls []as.Address
}

func NewStubAddressStore(data []as.Address) *StubAddressStore {
	return &StubAddressStore{
		addresses:   data,
		storeCalls:  []as.Address{},
		deleteCalls: []int{},
		updateCalls: []as.Address{},
	}
}

func (s *StubAddressStore) GetAddressesByCustomerId(customerId int) ([]as.Address, error) {
	if customerId == td.PeterCustomer.Id {
		return []as.Address{td.PeterAddress1, td.PeterAddress2}, nil
	}

	if customerId == td.AliceCustomer.Id {
		return []as.Address{td.AliceAddress}, nil
	}

	return []as.Address{}, nil
}

func (s *StubAddressStore) StoreAddress(address as.Address) {
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

func (s *StubAddressStore) GetAddressById(id int) (as.Address, error) {
	for _, address := range s.addresses {
		if address.Id == id {
			return address, nil
		}
	}
	return as.Address{}, fmt.Errorf("address with id %d doesn't exist", id)
}

func (s *StubAddressStore) UpdateAddress(address as.Address) error {
	s.updateCalls = append(s.updateCalls, address)
	return nil
}

func (s *StubAddressStore) Empty() {
	s.addresses = []as.Address{}
	s.storeCalls = []as.Address{}
}
