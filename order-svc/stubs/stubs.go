package stubs

import (
	"strconv"

	"github.com/VitoNaychev/food-app/msgtypes"
	"github.com/VitoNaychev/food-app/order-svc/models"
	"github.com/VitoNaychev/food-app/storeerrors"
)

func StubVerifyJWT(jwt string) (msgtypes.AuthResponse, error) {
	if jwt == "invalidJWT" {
		return msgtypes.AuthResponse{Status: msgtypes.INVALID, ID: 0}, nil
	} else if jwt == "10" {
		return msgtypes.AuthResponse{Status: msgtypes.NOT_FOUND, ID: 0}, nil
	} else {
		id, _ := strconv.Atoi(jwt)
		return msgtypes.AuthResponse{Status: msgtypes.OK, ID: id}, nil
	}
}

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
	return models.Address{}, storeerrors.ErrNotFound
}

type StubOrderStore struct {
	CreatedOrders []models.Order
	Orders        []models.Order
}

func (s *StubOrderStore) GetOrderByID(id int) (models.Order, error) {
	for _, order := range s.Orders {
		if order.ID == id {
			return order, nil
		}
	}
	return models.Order{}, storeerrors.ErrNotFound
}

func (s *StubOrderStore) CancelOrder(id int) error {
	for i := range s.Orders {
		if s.Orders[i].ID == id {
			s.Orders[i].Status = models.CANCELED
			return nil
		}
	}
	return nil
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
		if order.CustomerID == customerID &&
			order.Status != models.CANCELED && order.Status != models.COMPLETED {
			customerOrders = append(customerOrders, order)
		}
	}
	return customerOrders, nil
}
