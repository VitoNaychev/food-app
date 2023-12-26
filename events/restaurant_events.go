package events

import "github.com/VitoNaychev/food-app/restaurant-svc/models"

const RESTAURANT_EVENTS_TOPIC = "restaurant-events-topic"

type MenuItemCreatedEvent struct {
	ID           int     `validate:"min=1"             json:"id"`
	RestaurantID int     `validate:"min=1"             json:"restaurant_id"`
	Name         string  `validate:"min=2,max=20"      json:"name"`
	Price        float32 `validate:"required,max=1000" json:"price"`
}

func NewMenuItemCreatedEvent(menuItem models.MenuItem, restaurantID int) MenuItemCreatedEvent {
	event := MenuItemCreatedEvent{
		ID:           menuItem.ID,
		RestaurantID: restaurantID,
		Name:         menuItem.Name,
		Price:        menuItem.Price,
	}

	return event
}
