package handlers_test

import (
	"testing"

	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/kitchen-svc/handlers"
	"github.com/VitoNaychev/food-app/kitchen-svc/models"
	"github.com/VitoNaychev/food-app/kitchen-svc/testdata"
	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/VitoNaychev/food-app/testutil"
)

type StubMenuItemStore struct {
	menuItems                []models.MenuItem
	createdMenuItem          models.MenuItem
	deletedMenuItemID        int
	updatedMenuItem          models.MenuItem
	deletedItemsRestaurantID int
}

func (s *StubMenuItemStore) GetMenuItemByID(id int) (models.MenuItem, error) {
	for _, menuItem := range s.menuItems {
		if menuItem.ID == id {
			return menuItem, nil
		}
	}
	return models.MenuItem{}, storeerrors.ErrNotFound
}

func (s *StubMenuItemStore) DeleteMenuItemWhereRestaurantID(id int) error {
	s.deletedItemsRestaurantID = id
	return nil
}

func (s *StubMenuItemStore) CreateMenuItem(menuItem *models.MenuItem) error {
	s.createdMenuItem = *menuItem
	return nil
}

func (s *StubMenuItemStore) DeleteMenuItem(id int) error {
	s.deletedMenuItemID = id
	return nil
}

func (s *StubMenuItemStore) UpdateMenuItem(menuItem *models.MenuItem) error {
	s.updatedMenuItem = *menuItem
	return nil
}

type StubRestaurantStore struct {
	restaurants         []models.Restaurant
	createdRestaurant   models.Restaurant
	deletedRestaurantID int
}

func (s *StubRestaurantStore) DeleteRestaurant(id int) error {
	s.deletedRestaurantID = id
	return nil
}

func (s *StubRestaurantStore) GetRestaurantByID(id int) (models.Restaurant, error) {
	for _, restaurant := range s.restaurants {
		if restaurant.ID == id {
			return restaurant, nil
		}
	}

	return models.Restaurant{}, nil
}

func (s *StubRestaurantStore) CreateRestaurant(restaurant *models.Restaurant) error {
	s.createdRestaurant = *restaurant

	return nil
}

func TestRestaurantEventHandler(t *testing.T) {
	restaurantStore := &StubRestaurantStore{}
	menuItemStore := &StubMenuItemStore{}
	restaurantEventHandler := handlers.NewRestaurantEventHandler(restaurantStore, menuItemStore)

	t.Run("creates restaurant on RESTAURANT_CREATED_EVENT", func(t *testing.T) {
		payload := events.RestaurantCreatedEvent{ID: testdata.ShackRestaurant.ID}
		event := events.NewTypedEvent(events.RESTAURANT_CREATED_EVENT_ID, testdata.ShackRestaurant.ID, payload)

		err := restaurantEventHandler.HandleRestaurantCreatedEvent(event)
		testutil.AssertNoErr(t, err)

		got := restaurantStore.createdRestaurant
		testutil.AssertEqual(t, got, testdata.ShackRestaurant)
	})

	t.Run("deletes restaurant on RESTAURANT_DELETED_EVENT and all related menu items", func(t *testing.T) {
		payload := events.RestaurantDeletedEvent{ID: testdata.ShackRestaurant.ID}
		event := events.NewTypedEvent(events.RESTAURANT_DELETED_EVENT_ID, testdata.ShackRestaurant.ID, payload)

		err := restaurantEventHandler.HandleRestaurantDeletedEvent(event)
		testutil.AssertNoErr(t, err)

		testutil.AssertEqual(t, restaurantStore.deletedRestaurantID, testdata.ShackRestaurant.ID)
		testutil.AssertEqual(t, menuItemStore.deletedItemsRestaurantID, testdata.ShackRestaurant.ID)
	})
}

func TestRestaurantMenuEventHandler(t *testing.T) {
	menuItemStore := &StubMenuItemStore{}
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

		got := menuItemStore.createdMenuItem
		testutil.AssertEqual(t, got, testdata.ShackMenuItem)
	})

	t.Run("deletes menu item on MENU_ITEM_DELETED_EVENT", func(t *testing.T) {
		payload := events.MenuItemDeletedEvent{ID: testdata.ShackMenuItem.ID}
		event := events.NewTypedEvent(events.MENU_ITEM_DELETED_EVENT_ID, testdata.ShackRestaurant.ID, payload)

		err := restaurantEventHandler.HandleMenuItemDeletedEvent(event)
		testutil.AssertNoErr(t, err)

		got := menuItemStore.deletedMenuItemID
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

		got := menuItemStore.updatedMenuItem
		testutil.AssertEqual(t, got, testdata.ShackMenuItem)
	})
}
