package models

import (
	"github.com/VitoNaychev/food-app/storeerrors"
)

type InMemoryOrderStore struct {
	orders []Order
}

func NewInMemoryOrderStore() *InMemoryOrderStore {
	return &InMemoryOrderStore{[]Order{}}
}

func (i *InMemoryOrderStore) GetOrderByID(id int) (Order, error) {
	for _, order := range i.orders {
		if order.ID == id {
			return order, nil
		}
	}

	return Order{}, storeerrors.ErrNotFound
}

func (i *InMemoryOrderStore) GetOrdersByCustomerID(customerID int) ([]Order, error) {
	var customerOrders []Order
	for _, order := range i.orders {
		if order.CustomerID == customerID {
			customerOrders = append(customerOrders, order)
		}
	}

	return customerOrders, nil
}

func (i *InMemoryOrderStore) GetCurrentOrdersByCustomerID(customerID int) ([]Order, error) {
	var currentOrders []Order
	for _, order := range i.orders {
		if order.CustomerID == customerID &&
			order.Status != CANCELED && order.Status != COMPLETED {
			currentOrders = append(currentOrders, order)
		}
	}

	return currentOrders, nil
}

func (i *InMemoryOrderStore) CreateOrder(order *Order) error {
	order.ID = len(i.orders) + 1
	i.orders = append(i.orders, *order)

	return nil
}

func (i *InMemoryOrderStore) CancelOrder(id int) error {
	for j, order := range i.orders {
		if order.ID == id {
			i.orders[j].Status = CANCELED
			return nil
		}
	}

	return storeerrors.ErrNotFound
}
