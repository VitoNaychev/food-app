package stubs

import (
	"github.com/VitoNaychev/food-app/order-svc/models"
	"github.com/VitoNaychev/food-app/storeerrors"
)

type StubOrderItemStore struct {
	CreatedOrderItems []models.OrderItem
	OrderItems        []models.OrderItem
}

func (s *StubOrderItemStore) CreateOrderItem(orderItem *models.OrderItem) error {
	orderItem.ID = len(s.CreatedOrderItems) + 1
	s.CreatedOrderItems = append(s.CreatedOrderItems, *orderItem)

	return nil
}

func (s *StubOrderItemStore) GetOrderItemsByOrderID(orderID int) ([]models.OrderItem, error) {
	orderItems := []models.OrderItem{}

	for _, orderItem := range s.OrderItems {
		if orderItem.OrderID == orderID {
			orderItems = append(orderItems, orderItem)
		}
	}

	if len(orderItems) == 0 {
		return nil, storeerrors.ErrNotFound
	}

	return orderItems, nil
}
