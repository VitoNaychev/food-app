package testutil

import (
	"testing"

	"github.com/VitoNaychev/bt-order-svc/models"
)

type ErroneousAddressStore struct {
	err error
}

func (s *ErroneousAddressStore) CreateAddress(address *models.Address) error {
	return s.err
}

func (s *ErroneousAddressStore) GetAddressByID(id int) (models.Address, error) {
	return models.Address{}, s.err
}

type ErroneousOrderStore struct {
	err error
}

func (s *ErroneousOrderStore) CreateOrder(order *models.Order) error {
	return s.err
}

func (s *ErroneousOrderStore) GetOrdersByCustomerID(customerID int) ([]models.Order, error) {
	return []models.Order{}, s.err
}

func (s *ErroneousOrderStore) GetCurrentOrdersByCustomerID(customerID int) ([]models.Order, error) {
	return []models.Order{}, s.err
}

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
	return models.Address{}, models.ErrNotFound
}

type StubOrderStore struct {
	CreatedOrders []models.Order
	Orders        []models.Order
}

func (s *StubOrderStore) CreateOrder(order *models.Order) error {
	s.CreatedOrders = append(s.CreatedOrders, *order)
	order.ID = len(s.CreatedOrders)

	return nil
}

func (s *StubOrderStore) GetOrdersByCustomerID(customerID int) ([]models.Order, error) {
	var customerOrders []models.Order
	for _, order := range s.Orders {
		if order.CustomerID == customerID {
			customerOrders = append(customerOrders, order)
		}
	}
	return customerOrders, nil
}

func (s *StubOrderStore) GetCurrentOrdersByCustomerID(customerID int) ([]models.Order, error) {
	var customerOrders []models.Order
	for _, order := range s.Orders {
		if order.CustomerID == customerID && order.Status != models.COMPLETED {
			customerOrders = append(customerOrders, order)
		}
	}
	return customerOrders, nil
}

func AssertStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Fatalf("got status %v want %v", got, want)
	}
}
