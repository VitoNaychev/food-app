package models

import (
	"github.com/VitoNaychev/food-app/storeerrors"
)

type InMemoryDeliveryStore struct {
	deliveries []Delivery
}

func NewInMemoryDeliveryStore() *InMemoryDeliveryStore {
	return &InMemoryDeliveryStore{[]Delivery{}}
}

func (i *InMemoryDeliveryStore) CreateDelivery(delivery *Delivery) error {
	i.deliveries = append(i.deliveries, *delivery)

	return nil
}

func (i *InMemoryDeliveryStore) GetDeliveryByID(id int) (Delivery, error) {
	for _, delivery := range i.deliveries {
		if delivery.ID == id {
			return delivery, nil
		}
	}
	return Delivery{}, storeerrors.ErrNotFound
}

func (i *InMemoryDeliveryStore) UpdateDelivery(updatedDelivery *Delivery) error {
	for j, delivery := range i.deliveries {
		if delivery.ID == updatedDelivery.ID {
			i.deliveries[j] = *updatedDelivery
			return nil
		}
	}
	return storeerrors.ErrNotFound
}

func (i *InMemoryDeliveryStore) GetActiveDeliveryByCourierID(courierID int) (Delivery, error) {
	for _, delivery := range i.deliveries {
		if delivery.CourierID == courierID &&
			delivery.State != CANCELED &&
			delivery.State != COMPLETED &&
			delivery.State != DECLINED {
			return delivery, nil
		}
	}
	return Delivery{}, storeerrors.ErrNotFound
}
