package events

import "github.com/VitoNaychev/food-app/restaurant-svc/models"

const RESTAURANT_EVENTS_TOPIC = "restaurant-events-topic"

const (
	RESTAURANT_CREATED_EVENT_ID EventID = iota
	RESTAURANT_DELETED_EVENT_ID
	MENU_ITEM_CREATED_EVENT_ID
	MENU_ITEM_DELETED_EVENT_ID
	MENU_ITEM_UPDATED_EVENT_ID
)

type RestaurantCreatedEvent struct {
	ID int `validate:"min=1"             json:"id"`
}

type RestaurantDeletedEvent struct {
	ID int `validate:"min=1"             json:"id"`
}

type MenuItemUpdatedEvent struct {
	ID           int     `validate:"min=1"             json:"id"`
	RestaurantID int     `validate:"min=1"             json:"restaurant_id"`
	Name         string  `validate:"min=2,max=20"      json:"name"`
	Price        float32 `validate:"required,max=1000" json:"price"`
}

func NewMenuItemUpdatedEvent(menuItem models.MenuItem) MenuItemUpdatedEvent {
	event := MenuItemUpdatedEvent{
		ID:           menuItem.ID,
		RestaurantID: menuItem.RestaurantID,
		Name:         menuItem.Name,
		Price:        menuItem.Price,
	}

	return event
}

type MenuItemDeletedEvent struct {
	ID int `validate:"min=1"             json:"id"`
}

type MenuItemCreatedEvent struct {
	ID           int     `validate:"min=1"             json:"id"`
	RestaurantID int     `validate:"min=1"             json:"restaurant_id"`
	Name         string  `validate:"min=2,max=20"      json:"name"`
	Price        float32 `validate:"required,max=1000" json:"price"`
}

func NewMenuItemCreatedEvent(menuItem models.MenuItem) MenuItemCreatedEvent {
	event := MenuItemCreatedEvent{
		ID:           menuItem.ID,
		RestaurantID: menuItem.RestaurantID,
		Name:         menuItem.Name,
		Price:        menuItem.Price,
	}

	return event
}
