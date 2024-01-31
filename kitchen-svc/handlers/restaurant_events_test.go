package handlers_test

import (
	"testing"

	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/kitchen-svc/handlers"
	"github.com/VitoNaychev/food-app/kitchen-svc/stubs"
	"github.com/VitoNaychev/food-app/kitchen-svc/testdata"

	"github.com/VitoNaychev/food-app/testutil"
)

func TestRestaurantEventHandler(t *testing.T) {
	restaurantStore := &stubs.StubRestaurantStore{}
	menuItemStore := &stubs.StubMenuItemStore{}
	restaurantEventHandler := handlers.NewRestaurantEventHandler(restaurantStore, menuItemStore)

	t.Run("creates restaurant on RESTAURANT_CREATED_EVENT", func(t *testing.T) {
		payload := events.RestaurantCreatedEvent{ID: testdata.ShackRestaurant.ID}
		event := events.NewTypedEvent(events.RESTAURANT_CREATED_EVENT_ID, testdata.ShackRestaurant.ID, payload)

		err := restaurantEventHandler.HandleRestaurantCreatedEvent(event)
		testutil.AssertNoErr(t, err)

		got := restaurantStore.CreatedRestaurant
		testutil.AssertEqual(t, got, testdata.ShackRestaurant)
	})

	t.Run("deletes restaurant on RESTAURANT_DELETED_EVENT and all related menu items", func(t *testing.T) {
		payload := events.RestaurantDeletedEvent{ID: testdata.ShackRestaurant.ID}
		event := events.NewTypedEvent(events.RESTAURANT_DELETED_EVENT_ID, testdata.ShackRestaurant.ID, payload)

		err := restaurantEventHandler.HandleRestaurantDeletedEvent(event)
		testutil.AssertNoErr(t, err)

		testutil.AssertEqual(t, restaurantStore.DeletedRestaurantID, testdata.ShackRestaurant.ID)
		testutil.AssertEqual(t, menuItemStore.DeletedItemsRestaurantID, testdata.ShackRestaurant.ID)
	})
}

func TestRestaurantMenuEventHandler(t *testing.T) {
	menuItemStore := &stubs.StubMenuItemStore{}
	restaurantEventHandler := handlers.NewRestaurantEventHandler(nil, menuItemStore)

	t.Run("creates menu item on MENU_ITEM_CREATED_EVENT", func(t *testing.T) {
		payload := events.MenuItemCreatedEvent{
			ID:           testdata.ShackMenuItem.ID,
			RestaurantID: testdata.ShackMenuItem.RestaurantID,
			Name:         testdata.ShackMenuItem.Name,
			Price:        testdata.ShackMenuItem.Price,
		}
		event := events.NewTypedEvent(events.MENU_ITEM_CREATED_EVENT_ID, testdata.ShackRestaurant.ID, payload)

		err := restaurantEventHandler.HandleMenuItemCreatedEvent(event)
		testutil.AssertNoErr(t, err)

		got := menuItemStore.CreatedMenuItem
		testutil.AssertEqual(t, got, testdata.ShackMenuItem)
	})

	t.Run("deletes menu item on MENU_ITEM_DELETED_EVENT", func(t *testing.T) {
		payload := events.MenuItemDeletedEvent{ID: testdata.ShackMenuItem.ID}
		event := events.NewTypedEvent(events.MENU_ITEM_DELETED_EVENT_ID, testdata.ShackRestaurant.ID, payload)

		err := restaurantEventHandler.HandleMenuItemDeletedEvent(event)
		testutil.AssertNoErr(t, err)

		got := menuItemStore.DeletedMenuItemID
		testutil.AssertEqual(t, got, testdata.ShackMenuItem.ID)
	})

	t.Run("updates menu item on MENU_ITEM_UPDATED_EVENT", func(t *testing.T) {
		payload := events.MenuItemUpdatedEvent{
			ID:           testdata.ShackMenuItem.ID,
			RestaurantID: testdata.ShackMenuItem.RestaurantID,
			Name:         testdata.ShackMenuItem.Name,
			Price:        testdata.ShackMenuItem.Price,
		}
		event := events.NewTypedEvent(events.MENU_ITEM_UPDATED_EVENT_ID, testdata.ShackRestaurant.ID, payload)

		err := restaurantEventHandler.HandleMenuItemUpdatedEvent(event)
		testutil.AssertNoErr(t, err)

		got := menuItemStore.UpdatedMenuItem
		testutil.AssertEqual(t, got, testdata.ShackMenuItem)
	})
}
