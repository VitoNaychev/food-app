package models

type OrderItemStore interface {
	CreateOrderItem(*OrderItem) error
	GetOrderItemsByOrderID(int) ([]OrderItem, error)
}
