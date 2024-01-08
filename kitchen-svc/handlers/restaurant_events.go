package handlers

import (
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

func (r *RestaurantEventHandler) HandleRestaurantCreatedEvent(event events.Event[events.RestaurantCreatedEvent]) error {
	restaurant := models.Restaurant{ID: event.Payload.ID}
	r.restaurantStore.CreateRestaurant(&restaurant)
	return nil
}
