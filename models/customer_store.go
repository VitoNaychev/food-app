package models

type CustomerStore interface {
	GetCustomerByID(id int) (Customer, error)
	GetCustomerByEmail(email string) (Customer, error)
	CreateCustomer(customer *Customer) error
	DeleteCustomer(id int) error
	UpdateCustomer(customer *Customer) error
}
