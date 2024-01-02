package models

import "github.com/VitoNaychev/food-app/storeerrors"

type InMemoryAddressStore struct {
	addresses []Address
}

func NewInMemoryAddressStore() *InMemoryAddressStore {
	return &InMemoryAddressStore{[]Address{}}
}

func (i *InMemoryAddressStore) CreateAddress(address *Address) error {
	address.ID = len(i.addresses) + 1
	i.addresses = append(i.addresses, *address)
	return nil
}

func (i *InMemoryAddressStore) GetAddressByID(id int) (Address, error) {
	for _, address := range i.addresses {
		if address.ID == id {
			return address, nil
		}
	}
	return Address{}, storeerrors.ErrNotFound
}

func (i *InMemoryAddressStore) GetAddressByRestaurantID(restaurantID int) (Address, error) {
	for _, address := range i.addresses {
		if address.RestaurantID == restaurantID {
			return address, nil
		}
	}
	return Address{}, storeerrors.ErrNotFound
}

func (i *InMemoryAddressStore) UpdateAddress(address *Address) error {
	for j, oldAddress := range i.addresses {
		if oldAddress.ID == address.ID {
			i.addresses[j] = *address
			return nil
		}
	}
	return storeerrors.ErrNotFound
}
