package services_test

import (
	"testing"

	"github.com/VitoNaychev/food-app/kitchen-svc/domain"
	"github.com/VitoNaychev/food-app/kitchen-svc/services"
	"github.com/VitoNaychev/food-app/testutil"
)

var restaurant = domain.Restaurant{
	ID: 1,
}

type StubRestaurantStore struct {
	restaurant domain.Restaurant
}

func (s *StubRestaurantStore) CreateRestaurant(restaurant *domain.Restaurant) error {
	s.restaurant = *restaurant
	return nil
}

func (s *StubRestaurantStore) GetRestaurantByID(id int) (domain.Restaurant, error) {
	return s.restaurant, nil
}

func (s *StubRestaurantStore) UpdateRestaurant(restaurant *domain.Restaurant) error {
	panic("unimplemented")
}

func TestKithenService(t *testing.T) {
	store := &StubRestaurantStore{}
	service := services.NewKitchenService(store)

	t.Run("creates a restaurant", func(t *testing.T) {
		service.CreateRestaurant(restaurant.ID)

		testutil.AssertEqual(t, store.restaurant, restaurant)
	})

	t.Run("gets restaurant by id", func(t *testing.T) {
		got, _ := service.GetRestaurantByID(restaurant.ID)

		testutil.AssertEqual(t, got, restaurant)
	})
}
