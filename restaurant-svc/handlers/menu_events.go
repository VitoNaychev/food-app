package handlers

import (
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

func NewMenuItemCreatedEvent(menuItem models.MenuItem) events.MenuItemCreatedEvent {
	event := events.MenuItemCreatedEvent{
		ID:           menuItem.ID,
		RestaurantID: menuItem.RestaurantID,
		Name:         menuItem.Name,
		Price:        menuItem.Price,
	}

	return event
}

func NewMenuItemUpdatedEvent(menuItem models.MenuItem) events.MenuItemUpdatedEvent {
	event := events.MenuItemUpdatedEvent{
		ID:           menuItem.ID,
		RestaurantID: menuItem.RestaurantID,
		Name:         menuItem.Name,
		Price:        menuItem.Price,
	}

	return event
}
