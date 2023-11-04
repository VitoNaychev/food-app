package testutil

import (
	"github.com/VitoNaychev/bt-customer-svc/models"
)

type StubCustomerStore struct {
	customers   []models.Customer
	storeCalls  []models.Customer
	deleteCalls []int
	updateCalls []models.Customer
}

func NewStubCustomerStore(data []models.Customer) *StubCustomerStore {
	return &StubCustomerStore{
		customers:   data,
		storeCalls:  []models.Customer{},
		deleteCalls: []int{},
		updateCalls: []models.Customer{},
	}
}

func (s *StubCustomerStore) GetCustomerByID(id int) (models.Customer, error) {
	for _, customer := range s.customers {
		if customer.Id == id {
			return customer, nil
		}
	}

	return models.Customer{}, models.ErrNotFound
}

func (s *StubCustomerStore) GetCustomerByEmail(email string) (models.Customer, error) {
	for _, customer := range s.customers {
		if customer.Email == email {
			return customer, nil
		}
	}

	return models.Customer{}, models.ErrNotFound
}

func (s *StubCustomerStore) CreateCustomer(customer *models.Customer) error {
	customer.Id = len(s.customers) + 1
	s.customers = append(s.customers, *customer)
	s.storeCalls = append(s.storeCalls, *customer)

	return nil
}

func (s *StubCustomerStore) UpdateCustomer(customer *models.Customer) error {
	s.updateCalls = append(s.updateCalls, *customer)

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

	return models.ErrNotFound
}

func (s *StubCustomerStore) Empty() {
	s.customers = []models.Customer{}
	s.storeCalls = []models.Customer{}
}
