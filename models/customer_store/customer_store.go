package customer_store

type CustomerStore interface {
	GetCustomerById(id int) (Customer, error)
	GetCustomerByEmail(email string) (Customer, error)
	StoreCustomer(customer Customer) int
	DeleteCustomer(id int) error
	UpdateCustomer(customer Customer) error
}
