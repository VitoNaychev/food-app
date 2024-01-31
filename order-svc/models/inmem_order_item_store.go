package models

import (
	"github.com/VitoNaychev/food-app/storeerrors"
)

type InMemoryOrderItemStore struct {
	orderItems []OrderItem
}

func NewInMemoryOrderItemStore() *InMemoryOrderItemStore {
	return &InMemoryOrderItemStore{[]OrderItem{}}
}

func (i *InMemoryOrderItemStore) CreateOrderItem(orderItem *OrderItem) error {
	orderItem.ID = len(i.orderItems) + 1
	i.orderItems = append(i.orderItems, *orderItem)

	return nil
}

func (i *InMemoryOrderItemStore) GetOrderItemsByOrderID(orderID int) ([]OrderItem, error) {
	var orderItems []OrderItem
	for _, item := range i.orderItems {
		if item.OrderID == orderID {
			orderItems = append(orderItems, item)
		}
	}

	if len(orderItems) == 0 {
		return nil, storeerrors.ErrNotFound
	}

	return orderItems, nil
}
