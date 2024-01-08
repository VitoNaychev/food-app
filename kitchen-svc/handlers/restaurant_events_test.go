package handlers_test

import (
	"testing"
	"time"

	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/kitchen-svc/handlers"
	"github.com/VitoNaychev/food-app/kitchen-svc/models"
	"github.com/VitoNaychev/food-app/kitchen-svc/testdata"
	"github.com/VitoNaychev/food-app/testutil"
)

type StubRestaurantStore struct {
	createdRestaurant models.Restaurant
}

func (*StubRestaurantStore) DeleteRestaurant(int) error {
	panic("unimplemented")
}

func (*StubRestaurantStore) GetRestaurantByID(int) (models.Restaurant, error) {
	panic("unimplemented")
}

func (s *StubRestaurantStore) CreateRestaurant(restaurant *models.Restaurant) error {
	s.createdRestaurant = *restaurant

	return nil
}

func TestRestaurantEventHandler(t *testing.T) {
	restaurantStore := &StubRestaurantStore{}
	restaurantEventHandler := handlers.NewRestaurantEventHandler(restaurantStore)

	t.Run("creates restaurant on RESTAURANT_CREATED_EVENT", func(t *testing.T) {
		payload := events.RestaurantCreatedEvent{ID: testdata.ShackRestaurant.ID}
		event := events.Event[events.RestaurantCreatedEvent]{
			EventID:     events.RESTAURANT_CREATED_EVENT_ID,
			AggregateID: testdata.ShackRestaurant.ID,
			Timestamp:   time.Now().Round(0),
			Payload:     payload,
		}

		err := restaurantEventHandler.HandleRestaurantCreatedEvent(event)
		testutil.AssertNoErr(t, err)

		got := restaurantStore.createdRestaurant
		testutil.AssertEqual(t, got, testdata.ShackRestaurant)
	})
}
