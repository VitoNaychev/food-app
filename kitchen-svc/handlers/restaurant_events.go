package handlers

import (
	"reflect"

	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/kitchen-svc/models"
)

type RestaurantEventHandler struct {
	restaurantStore models.RestaurantStore
	menuItemStore   models.MenuItemStore
}

func NewRestaurantEventHandler(restaurantStore models.RestaurantStore, menuItemStore models.MenuItemStore) *RestaurantEventHandler {
	endpoint := RestaurantEventHandler{
		restaurantStore: restaurantStore,
		menuItemStore:   menuItemStore,
	}

	return &endpoint
}

func RegisterRestaurantEventHandlers(eventConsumer events.EventConsumer, restaurantEventHandler *RestaurantEventHandler) {
	eventConsumer.RegisterEventHandler(events.RESTAURANT_EVENTS_TOPIC,
		events.RESTAURANT_CREATED_EVENT_ID,
		events.EventHandlerWrapper(restaurantEventHandler.HandleRestaurantCreatedEvent),
		reflect.TypeOf(events.RestaurantCreatedEvent{}))
	eventConsumer.RegisterEventHandler(events.RESTAURANT_EVENTS_TOPIC,
		events.RESTAURANT_DELETED_EVENT_ID,
		events.EventHandlerWrapper(restaurantEventHandler.HandleRestaurantDeletedEvent),
		reflect.TypeOf(events.RestaurantDeletedEvent{}))
	eventConsumer.RegisterEventHandler(events.RESTAURANT_EVENTS_TOPIC,
		events.MENU_ITEM_CREATED_EVENT_ID,
		events.EventHandlerWrapper(restaurantEventHandler.HandleMenuItemCreatedEvent),
		reflect.TypeOf(events.MenuItemCreatedEvent{}))
	eventConsumer.RegisterEventHandler(events.RESTAURANT_EVENTS_TOPIC,
		events.MENU_ITEM_UPDATED_EVENT_ID,
		events.EventHandlerWrapper(restaurantEventHandler.HandleMenuItemUpdatedEvent),
		reflect.TypeOf(events.MenuItemUpdatedEvent{}))
	eventConsumer.RegisterEventHandler(events.RESTAURANT_EVENTS_TOPIC,
		events.MENU_ITEM_DELETED_EVENT_ID,
		events.EventHandlerWrapper(restaurantEventHandler.HandleMenuItemDeletedEvent),
		reflect.TypeOf(events.MenuItemDeletedEvent{}))
}

func (r *RestaurantEventHandler) HandleRestaurantCreatedEvent(event events.Event[events.RestaurantCreatedEvent]) error {
	restaurant := models.Restaurant{ID: event.Payload.ID}
	err := r.restaurantStore.CreateRestaurant(&restaurant)
	return err
}

func (r *RestaurantEventHandler) HandleRestaurantDeletedEvent(event events.Event[events.RestaurantDeletedEvent]) error {
	err := r.restaurantStore.DeleteRestaurant(event.Payload.ID)
	if err != nil {
		return err
	}

	err = r.menuItemStore.DeleteMenuItemWhereRestaurantID(event.Payload.ID)
	return err
}

func (r *RestaurantEventHandler) HandleMenuItemCreatedEvent(event events.Event[events.MenuItemCreatedEvent]) error {
	menuItem := models.MenuItem{
		ID:           event.Payload.ID,
		RestaurantID: event.Payload.RestaurantID,
		Name:         event.Payload.Name,
		Price:        event.Payload.Price,
	}
	err := r.menuItemStore.CreateMenuItem(&menuItem)
	return err
}

func (r *RestaurantEventHandler) HandleMenuItemDeletedEvent(event events.Event[events.MenuItemDeletedEvent]) error {
	err := r.menuItemStore.DeleteMenuItem(event.Payload.ID)
	return err
}

func (r *RestaurantEventHandler) HandleMenuItemUpdatedEvent(event events.Event[events.MenuItemUpdatedEvent]) error {
	menuItem := models.MenuItem{
		ID:           event.Payload.ID,
		RestaurantID: event.Payload.RestaurantID,
		Name:         event.Payload.Name,
		Price:        event.Payload.Price,
	}
	err := r.menuItemStore.UpdateMenuItem(&menuItem)
	return err
}
