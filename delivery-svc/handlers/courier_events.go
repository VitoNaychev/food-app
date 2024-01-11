package handlers

import (
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/events/svcevents"
)

type CourierEventHandler struct {
	courierStore models.CourierStore
}

func NewCourierEventHandler(courierStore models.CourierStore) *CourierEventHandler {
	courierEventHandler := CourierEventHandler{
		courierStore: courierStore,
	}

	return &courierEventHandler
}

func (c *CourierEventHandler) HandleCourierCreatedEvent(event events.Event[svcevents.CourierCreatedEvent]) error {
	courier := models.Courier{
		ID:   event.Payload.ID,
		Name: event.Payload.Name,
	}

	err := c.courierStore.CreateCourier(&courier)
	return err
}

func (c *CourierEventHandler) HandleCourierDeletedEvent(event events.Event[svcevents.CourierDeletedEvent]) error {
	err := c.courierStore.DeleteCourier(event.Payload.ID)
	return err
}
