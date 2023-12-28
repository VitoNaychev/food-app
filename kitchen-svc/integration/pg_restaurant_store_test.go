package integration

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/integrationutil"
	"github.com/VitoNaychev/food-app/kitchen-svc/domain"
	"github.com/VitoNaychev/food-app/kitchen-svc/handlers"
	"github.com/VitoNaychev/food-app/kitchen-svc/services"
	"github.com/VitoNaychev/food-app/kitchen-svc/stores"
	"github.com/VitoNaychev/food-app/testutil"
)

var restaurant = domain.Restaurant{
	ID: 5,
}

func TestPGRestaurantStore(t *testing.T) {
	connStr := integrationutil.SetupDatabaseContainer(t, env)

	store, err := stores.NewPgRestaurantStore(context.Background(), connStr)
	testutil.AssertNil(t, err)

	service := services.NewKitchenService(&store)
	endpoint := handlers.NewRestaurantEventEndpoint(service)

	payload := events.RestaurantCreatedEvent{ID: restaurant.ID}
	payloadJSON, _ := json.Marshal(payload)
	envelope := events.NewEventEnvelope(events.RESTAURANT_CREATED_EVENT_ID, restaurant.ID)

	endpoint.HandleRestaurantEvent(envelope, payloadJSON)

	got, err := service.GetRestaurantByID(restaurant.ID)
	testutil.AssertNil(t, err)
	testutil.AssertEqual(t, got, restaurant)
}
