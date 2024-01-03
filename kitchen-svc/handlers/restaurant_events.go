package handlers

import (
	"encoding/json"

	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/kitchen-svc/models"
)

type RestaurantEventHandler struct {
	restaurantStore models.RestaurantStore
}

func NewRestaurantEventHandler(restaurantStore models.RestaurantStore) *RestaurantEventHandler {
	endpoint := RestaurantEventHandler{
		restaurantStore: restaurantStore,
	}

	return &endpoint
}

func (r *RestaurantEventHandler) HandleRestaurantEvent(envelope events.EventEnvelope, payload []byte) error {
	switch envelope.EventID {
	case events.RESTAURANT_CREATED_EVENT_ID:
		var event events.RestaurantCreatedEvent
		json.Unmarshal(payload, &event)

		return r.HandleRestaurantCreatedEvent(envelope, event)
	}

	return nil
}

func (r *RestaurantEventHandler) HandleRestaurantCreatedEvent(envelope events.EventEnvelope, event events.RestaurantCreatedEvent) error {
	restaurant := models.Restaurant{ID: event.ID}
	r.restaurantStore.CreateRestaurant(&restaurant)
	return nil
}
