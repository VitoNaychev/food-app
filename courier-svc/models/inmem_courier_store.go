package models

import (
	"github.com/VitoNaychev/food-app/storeerrors"
)

type InMemoryCourierStore struct {
	couriers []Courier
}

func NewInMemoryCourierStore() *InMemoryCourierStore {
	return &InMemoryCourierStore{[]Courier{}}
}

func (i *InMemoryCourierStore) DeleteCourier(id int) error {
	for j, courier := range i.couriers {
		if courier.ID == id {
			i.couriers = append(i.couriers[:j], i.couriers[j+1:]...)
			return nil
		}
	}
	return storeerrors.ErrNotFound
}

func (i *InMemoryCourierStore) UpdateCourier(updatedCourier *Courier) error {
	for j, courier := range i.couriers {
		if courier.ID == updatedCourier.ID {
			i.couriers[j] = *updatedCourier
			return nil
		}
	}
	return storeerrors.ErrNotFound
}

func (i *InMemoryCourierStore) CreateCourier(courier *Courier) error {
	courier.ID = len(i.couriers) + 1
	i.couriers = append(i.couriers, *courier)

	return nil
}

func (i *InMemoryCourierStore) GetCourierByID(id int) (Courier, error) {
	for _, courier := range i.couriers {
		if courier.ID == id {
			return courier, nil
		}
	}
	return Courier{}, storeerrors.ErrNotFound
}

func (i *InMemoryCourierStore) GetCourierByEmail(email string) (Courier, error) {
	for _, courier := range i.couriers {
		if courier.Email == email {
			return courier, nil
		}
	}
	return Courier{}, storeerrors.ErrNotFound
}
