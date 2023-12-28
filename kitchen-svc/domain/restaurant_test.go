package domain_test

import (
	"testing"

	"github.com/VitoNaychev/food-app/kitchen-svc/domain"
	"github.com/VitoNaychev/food-app/testutil"
)

func TestRestaurantOperations(t *testing.T) {
	t.Run("creates new restaurant", func(t *testing.T) {
		id := 10
		restaurant, err := domain.NewRestaurant(id)

		testutil.AssertNil(t, err)
		testutil.AssertEqual(t, restaurant.ID, id)
	})

	t.Run("returns ErrInvalidID on invalid ID", func(t *testing.T) {
		id := -1
		_, err := domain.NewRestaurant(id)

		testutil.AssertEqual(t, err, domain.ErrInvalidID)
	})
}
