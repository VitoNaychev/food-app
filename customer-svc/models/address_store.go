package models

type CustomerAddressStore interface {
	GetAddressByID(id int) (Address, error)
	GetAddressesByCustomerID(customerID int) ([]Address, error)
	CreateAddress(address *Address) error
	DeleteAddress(id int) error
	UpdateAddress(address *Address) error
}
