package models

type AddressStore interface {
	GetAddressByID(id int) (Address, error)
	// CreateAddress(address *Address) error
	// DeleteAddress(id int) error
	// UpdateAddress(address *Address) error
}
