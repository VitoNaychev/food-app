package models

type AddressStore interface {
	GetAddressByID(id int) (Address, error)
	GetAddressByRestaurantID(restaurantID int) (Address, error)
	CreateAddress(address *Address) error
	UpdateAddress(address *Address) error
}
