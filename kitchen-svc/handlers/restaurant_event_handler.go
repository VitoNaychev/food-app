package handlers

import (
	"encoding/json"

	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/kitchen-svc/services"
)

type RestaurantEventEndpoint struct {
	service services.KitchenServiceInterface
}

func NewRestaurantEventEndpoint(service services.KitchenServiceInterface) *RestaurantEventEndpoint {
	endpoint := RestaurantEventEndpoint{
		service: service,
	}

	return &endpoint
}

func (r *RestaurantEventEndpoint) HandleRestaurantEvent(envelope events.EventEnvelope, payload []byte) {
	switch envelope.EventID {
	case events.RESTAURANT_CREATED_EVENT_ID:
		var event events.RestaurantCreatedEvent
		json.Unmarshal(payload, &event)

		r.HandleRestaurantCreatedEvent(envelope, event)
	}
}

func (r *RestaurantEventEndpoint) HandleRestaurantCreatedEvent(envelope events.EventEnvelope, event events.RestaurantCreatedEvent) {
	r.service.CreateRestaurant(event.ID)
}

func (r *RestaurantEventEndpoint) HandleRestaurantDeletedEvent(event events.MenuItemCreatedEvent) {
}

func (r *RestaurantEventEndpoint) HandleMenuItemCreatedEvent(event events.MenuItemCreatedEvent) {
	r.service.CreateMenuItem(event.ID, event.RestaurantID, event.Name, event.Price)
}

func (r *RestaurantEventEndpoint) HandleMenuItemUpdatedEvent(event events.MenuItemCreatedEvent) {

}

func (r *RestaurantEventEndpoint) HandleMenuItemDeketedEvent(event events.MenuItemCreatedEvent) {

}
