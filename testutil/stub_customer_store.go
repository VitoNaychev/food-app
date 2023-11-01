package testutil

import (
	"fmt"

	cs "github.com/VitoNaychev/bt-customer-svc/models/customer_store"
)

type StubCustomerStore struct {
	customers   []cs.Customer
	storeCalls  []cs.Customer
	deleteCalls []int
	updateCalls []cs.Customer
}

func NewStubCustomerStore(data []cs.Customer) *StubCustomerStore {
	return &StubCustomerStore{
		customers:   data,
		storeCalls:  []cs.Customer{},
		deleteCalls: []int{},
		updateCalls: []cs.Customer{},
	}
}

func (s *StubCustomerStore) UpdateCustomer(customer cs.Customer) error {
	s.updateCalls = append(s.updateCalls, customer)

	return nil
}

func (s *StubCustomerStore) DeleteCustomer(id int) error {
	for _, customer := range s.customers {
		if customer.Id == id {
			// s.customers = append(s.customers[:id], s.customers[id+1:]...)
			s.deleteCalls = append(s.deleteCalls, id)
			return nil
		}
	}

	return fmt.Errorf("no customer with id %d", id)
}

func (s *StubCustomerStore) GetCustomerById(id int) (cs.Customer, error) {
	for _, customer := range s.customers {
		if customer.Id == id {
			return customer, nil
		}
	}

	return cs.Customer{}, fmt.Errorf("no customer with id %d", id)
}

func (s *StubCustomerStore) GetCustomerByEmail(email string) (cs.Customer, error) {
	for _, customer := range s.customers {
		if customer.Email == email {
			return customer, nil
		}
	}

	return cs.Customer{}, fmt.Errorf("no customer with email %v", email)
}

func (s *StubCustomerStore) StoreCustomer(customer cs.Customer) int {
	customer.Id = len(s.customers) + 1
	s.customers = append(s.customers, customer)
	s.storeCalls = append(s.storeCalls, customer)

	return customer.Id
}

func (s *StubCustomerStore) Empty() {
	s.customers = []cs.Customer{}
	s.storeCalls = []cs.Customer{}
}
