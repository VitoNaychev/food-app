package models

import (
	"github.com/VitoNaychev/food-app/storeerrors"
)

type InMemoryLocationStore struct {
	locations []Location
}

func NewInMemoryLocationStore() *InMemoryLocationStore {
	return &InMemoryLocationStore{[]Location{}}
}

func (i *InMemoryLocationStore) CreateLocation(location *Location) error {
	i.locations = append(i.locations, *location)

	return nil
}

func (i *InMemoryLocationStore) GetLocationByCourierID(courierID int) (Location, error) {
	for _, location := range i.locations {
		if location.CourierID == courierID {
			return location, nil
		}
	}

	return Location{}, storeerrors.ErrNotFound
}

func (i *InMemoryLocationStore) UpdateLocation(updatedLocation *Location) error {
	for j, location := range i.locations {
		if location.CourierID == updatedLocation.CourierID {
			i.locations[j] = *updatedLocation
			return nil
		}
	}

	return storeerrors.ErrNotFound
}

func (i *InMemoryLocationStore) DeleteLocation(courierID int) error {
	for j, location := range i.locations {
		if location.CourierID == courierID {
			i.locations = append(i.locations[:j], i.locations[j+1:]...)
			return nil
		}
	}

	return storeerrors.ErrNotFound
}
