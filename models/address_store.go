package models

type CustomerAddressStore interface {
	GetAddressesByCustomerId(customerId int) ([]Address, error)
	StoreAddress(address Address)
	DeleteAddressById(id int) error
	GetAddressById(id int) (Address, error)
	UpdateAddress(address Address) error
}
