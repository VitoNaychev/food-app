package handlers_test

import (
	"encoding/json"
	"testing"

	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/kitchen-svc/domain"
	"github.com/VitoNaychev/food-app/kitchen-svc/handlers"
	"github.com/VitoNaychev/food-app/testutil"
)

var restaurant = domain.Restaurant{
	ID: 5,
}

type DummyService struct {
	createRestaurantID int
}

func (d *DummyService) CreateRestaurant(id int) error {
	d.createRestaurantID = id
	return nil
}

func (d *DummyService) CreateMenuItem(id int, restaurantID int, name string, price float32) error {
	return nil
}

func TestRestaurantEventHandlers(t *testing.T) {
	service := &DummyService{}
	endpoint := handlers.NewRestaurantEventEndpoint(service)

	t.Run("handles restaurant created events", func(t *testing.T) {
		payload := events.RestaurantCreatedEvent{ID: restaurant.ID}
		payloadJSON, _ := json.Marshal(payload)
		envelope := events.NewEventEnvelope(events.RESTAURANT_CREATED_EVENT_ID, restaurant.ID)

		endpoint.HandleRestaurantEvent(envelope, payloadJSON)

		testutil.AssertEqual(t, service.createRestaurantID, restaurant.ID)
	})
}
