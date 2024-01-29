package handlers

import (
	"reflect"

	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/events/svcevents"
)

type CourierEventHandler struct {
	courierStore  models.CourierStore
	locationStore models.LocationStore
}

func NewCourierEventHandler(courierStore models.CourierStore, locationStore models.LocationStore) *CourierEventHandler {
	courierEventHandler := CourierEventHandler{
		courierStore:  courierStore,
		locationStore: locationStore,
	}

	return &courierEventHandler
}

func RegisterCourierEventHandlers(eventConsumer events.EventConsumer, courierEventHandler *CourierEventHandler) {
	eventConsumer.RegisterEventHandler(svcevents.COURIER_EVENTS_TOPIC,
		svcevents.COURIER_CREATED_EVENT_ID,
		events.EventHandlerWrapper(courierEventHandler.HandleCourierCreatedEvent),
		reflect.TypeOf(svcevents.CourierCreatedEvent{}))
	eventConsumer.RegisterEventHandler(svcevents.COURIER_EVENTS_TOPIC,
		svcevents.COURIER_DELETED_EVENT_ID,
		events.EventHandlerWrapper(courierEventHandler.HandleCourierDeletedEvent),
		reflect.TypeOf(svcevents.CourierDeletedEvent{}))
}

func (c *CourierEventHandler) HandleCourierCreatedEvent(event events.Event[svcevents.CourierCreatedEvent]) error {
	courier := models.Courier{
		ID:   event.Payload.ID,
		Name: event.Payload.Name,
	}
	err := c.courierStore.CreateCourier(&courier)
	if err != nil {
		return err
	}

	location := models.Location{
		CourierID: courier.ID,
		Lat:       0.0,
		Lon:       0.0,
	}
	err = c.locationStore.CreateLocation(&location)
	if err != nil {
		return err
	}

	return nil
}

func (c *CourierEventHandler) HandleCourierDeletedEvent(event events.Event[svcevents.CourierDeletedEvent]) error {
	err := c.locationStore.DeleteLocation(event.Payload.ID)
	if err != nil {
		return err
	}

	err = c.courierStore.DeleteCourier(event.Payload.ID)
	if err != nil {
		return err
	}

	return nil
}
