package svcevents

import "github.com/VitoNaychev/food-app/events"

const ORDER_EVENTS_TOPIC = "order-events-topic"

const (
	ORDER_CREATED_EVENT_ID events.EventID = iota
)

type OrderCreatedEvent struct {
	ID              int
	RestaurantID    int
	Items           []OrderCreatedEventItem
	Total           float64
	PickupAddress   OrderCreatedEventAddress
	DeliveryAddress OrderCreatedEventAddress
}

type OrderCreatedEventItem struct {
	ID         int
	MenuItemID int
	Quantity   int
}

type OrderCreatedEventAddress struct {
	ID           int
	Lat          float64
	Lon          float64
	AddressLine1 string
	AddressLine2 string
	City         string
	Country      string
}
