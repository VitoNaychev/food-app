package stubs

import (
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/storeerrors"
)

type StubDeliveryStore struct {
	Deliveries      []models.Delivery
	UpdatedDelivery models.Delivery
}

func (d *StubDeliveryStore) GetActiveDeliveryByCourierID(courierID int) (models.Delivery, error) {
	for _, delivery := range d.Deliveries {
		if delivery.CourierID == courierID &&
			delivery.State != models.CANCELED &&
			delivery.State != models.COMPLETED &&
			delivery.State != models.DECLINED {
			return delivery, nil
		}
	}

	return models.Delivery{}, storeerrors.ErrNotFound
}

func (d *StubDeliveryStore) GetDeliveryByID(id int) (models.Delivery, error) {
	for _, delivery := range d.Deliveries {
		if delivery.ID == id {
			return delivery, nil
		}
	}

	return models.Delivery{}, storeerrors.ErrNotFound
}

func (d *StubDeliveryStore) UpdateDelivery(delivery *models.Delivery) error {
	d.UpdatedDelivery = *delivery

	return nil
}
