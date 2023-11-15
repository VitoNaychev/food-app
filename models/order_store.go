package models

type OrderStore interface {
	// GetOrderByID(id int) (Order, error)
	GetOrdersByCustomerID(customerID int) ([]Order, error)
	GetCurrentOrdersByCustomerID(customerID int) ([]Order, error)
	// CreateOrder(order *Order) error
	// DeleteOrder(id int) error
}
