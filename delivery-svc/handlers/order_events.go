package handlers

import (
	"reflect"
	"time"

	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/events/svcevents"
)

type OrderEventHandler struct {
	deliveryStore models.DeliveryStore
	addressStore  models.AddressStore
}

func NewOrderEventHandler(deliveryStore models.DeliveryStore, addressStore models.AddressStore) *OrderEventHandler {
	return &OrderEventHandler{
		deliveryStore: deliveryStore,
		addressStore:  addressStore,
	}
}

func RegisterOrderEventHandlers(eventConsumer events.EventConsumer, orderEventhandler *OrderEventHandler) {
	eventConsumer.RegisterEventHandler(svcevents.ORDER_EVENTS_TOPIC,
		svcevents.ORDER_CREATED_EVENT_ID,
		events.EventHandlerWrapper(orderEventhandler.HandleOrderCreatedEvent),
		reflect.TypeOf(svcevents.OrderCreatedEvent{}))
}

func (o *OrderEventHandler) HandleOrderCreatedEvent(event events.Event[svcevents.OrderCreatedEvent]) error {
	pickupAddress := AddressFromOrderCreatedEventAddress(event.Payload.PickupAddress)
	err := o.addressStore.CreateAddress(&pickupAddress)
	if err != nil {
		return err
	}

	deliveryAddress := AddressFromOrderCreatedEventAddress(event.Payload.DeliveryAddress)
	err = o.addressStore.CreateAddress(&deliveryAddress)
	if err != nil {
		return err
	}

	delivery := DeliveryFromOrderCreatedEvent(event.Payload)
	// AssignDeliveryToAvailableCourier(delivery)
	delivery.CourierID = 1
	delivery.ReadyBy = models.ZeroTime

	err = o.deliveryStore.CreateDelivery(&delivery)
	if err != nil {
		return err
	}

	return nil
}

func DeliveryFromOrderCreatedEvent(orderCreatedEvent svcevents.OrderCreatedEvent) models.Delivery {
	delivery := models.Delivery{
		ID:                orderCreatedEvent.ID,
		PickupAddressID:   orderCreatedEvent.PickupAddress.ID,
		DeliveryAddressID: orderCreatedEvent.DeliveryAddress.ID,
		ReadyBy:           time.Time{},
		State:             models.PENDING,
	}

	return delivery
}

func AddressFromOrderCreatedEventAddress(orderCreatedEventAddress svcevents.OrderCreatedEventAddress) models.Address {
	address := models.Address{
		ID:           orderCreatedEventAddress.ID,
		Lat:          orderCreatedEventAddress.Lat,
		Lon:          orderCreatedEventAddress.Lon,
		AddressLine1: orderCreatedEventAddress.AddressLine1,
		AddressLine2: orderCreatedEventAddress.AddressLine2,
		City:         orderCreatedEventAddress.City,
		Country:      orderCreatedEventAddress.Country,
	}

	return address
}
