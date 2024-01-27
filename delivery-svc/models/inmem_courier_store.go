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

func (i *InMemoryCourierStore) CreateCourier(courier *Courier) error {
	i.couriers = append(i.couriers, *courier)
	return nil
}

func (i *InMemoryCourierStore) GetAllCouriers() []Courier {
	return i.couriers
}

func (i *InMemoryCourierStore) GetCourierByID(id int) (Courier, error) {
	for _, courier := range i.couriers {
		if courier.ID == id {
			return courier, nil
		}
	}
	return Courier{}, storeerrors.ErrNotFound
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
