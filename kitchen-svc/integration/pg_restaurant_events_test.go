package integration

import (
	"context"
	"testing"

	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/integrationutil"
	"github.com/VitoNaychev/food-app/kitchen-svc/handlers"
	"github.com/VitoNaychev/food-app/kitchen-svc/models"
	"github.com/VitoNaychev/food-app/kitchen-svc/testdata"
	"github.com/VitoNaychev/food-app/storeerrors"

	"github.com/VitoNaychev/food-app/pgconfig"
	"github.com/VitoNaychev/food-app/testutil"
)

func TestRestaurantEventHandlerIntegration(t *testing.T) {
	config := pgconfig.GetConfigFromEnv(env)
	integrationutil.SetupDatabaseContainer(t, &config, "../sql-scripts/init.sql")

	connStr := config.GetConnectionString()

	restaurantStore, err := models.NewPgRestaurantStore(context.Background(), connStr)
	testutil.AssertNoErr(t, err)

	menuItemStore, err := models.NewPgMenuItemStore(context.Background(), connStr)
	testutil.AssertNoErr(t, err)

	restaurantEventHandler := handlers.NewRestaurantEventHandler(restaurantStore, menuItemStore)

	t.Run("creates new restaurant", func(t *testing.T) {
		want := testdata.ShackRestaurant

		payload := events.RestaurantCreatedEvent{ID: want.ID}
		event := events.NewTypedEvent(events.RESTAURANT_CREATED_EVENT_ID, want.ID, payload)

		restaurantEventHandler.HandleRestaurantCreatedEvent(event)

		got, err := restaurantStore.GetRestaurantByID(want.ID)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, got, want)
	})

	t.Run("destroys restaurant", func(t *testing.T) {
		want := testdata.ShackRestaurant

		payload := events.RestaurantDeletedEvent{ID: want.ID}
		event := events.NewTypedEvent(events.RESTAURANT_DELETED_EVENT_ID, want.ID, payload)

		restaurantEventHandler.HandleRestaurantDeletedEvent(event)

		_, err := restaurantStore.GetRestaurantByID(want.ID)

		testutil.AssertError(t, err, storeerrors.ErrNotFound)
	})

	{
		payload := events.RestaurantCreatedEvent{ID: testdata.ShackRestaurant.ID}
		event := events.NewTypedEvent(events.RESTAURANT_CREATED_EVENT_ID, testdata.ShackRestaurant.ID, payload)

		restaurantEventHandler.HandleRestaurantCreatedEvent(event)
	}

	t.Run("creates menu item", func(t *testing.T) {
		want := testdata.ShackMenuItem

		payload := events.MenuItemCreatedEvent{
			ID:           want.ID,
			RestaurantID: want.RestaurantID,
			Name:         want.Name,
			Price:        want.Price,
		}
		event := events.NewTypedEvent(events.MENU_ITEM_CREATED_EVENT_ID, testdata.ShackRestaurant.ID, payload)

		restaurantEventHandler.HandleMenuItemCreatedEvent(event)

		got, err := menuItemStore.GetMenuItemByID(want.ID)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, got, want)
	})

	t.Run("updates menu item", func(t *testing.T) {
		want := testdata.ShackMenuItem
		want.Name = "XXXXXL Duner"

		payload := events.MenuItemUpdatedEvent{
			ID:           want.ID,
			RestaurantID: want.RestaurantID,
			Name:         want.Name,
			Price:        want.Price,
		}
		event := events.NewTypedEvent(events.MENU_ITEM_UPDATED_EVENT_ID, testdata.ShackRestaurant.ID, payload)

		restaurantEventHandler.HandleMenuItemUpdatedEvent(event)

		got, err := menuItemStore.GetMenuItemByID(want.ID)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, got, want)
	})

	t.Run("deletes menu item", func(t *testing.T) {
		want := testdata.ShackMenuItem

		payload := events.MenuItemDeletedEvent{ID: want.ID}
		event := events.NewTypedEvent(events.MENU_ITEM_DELETED_EVENT_ID, testdata.ShackRestaurant.ID, payload)

		restaurantEventHandler.HandleMenuItemDeletedEvent(event)

		_, err := menuItemStore.GetMenuItemByID(want.ID)

		testutil.AssertError(t, err, storeerrors.ErrNotFound)
	})
}
