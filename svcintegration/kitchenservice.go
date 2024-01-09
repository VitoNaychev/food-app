package svcintegration

import (
	"context"
	"reflect"
	"testing"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/kitchen-svc/handlers"
	"github.com/VitoNaychev/food-app/kitchen-svc/models"
)

type KitchenService struct {
	restaurantStore models.RestaurantStore
	menuItemStore   models.MenuItemStore

	restaurantEventHandler *handlers.RestaurantEventHandler

	eventConsumer       *events.KafkaEventConsumer
	eventConsumerCtx    context.Context
	eventConsumerCancel context.CancelFunc
}

func SetupKitchenService(t testing.TB, env appenv.Enviornment, port string) KitchenService {
	eventConsumer, err := events.NewKafkaEventConsumer(env.KafkaBrokers, "kitchen-grp")
	if err != nil {
		t.Fatalf("Kafka Event Consumer error: %v\n", err)
	}

	restaurantStore := models.NewInMemoryRestaurantStore()
	menuItemStore := models.NewInMemoryMenuItemStore()
	restaurantEventHandler := handlers.NewRestaurantEventHandler(restaurantStore, menuItemStore)

	eventConsumerCtx, eventConsumerCancel := context.WithCancel(context.Background())

	kitchenService := KitchenService{
		restaurantStore: restaurantStore,

		restaurantEventHandler: restaurantEventHandler,
		menuItemStore:          menuItemStore,

		eventConsumer:       eventConsumer,
		eventConsumerCtx:    eventConsumerCtx,
		eventConsumerCancel: eventConsumerCancel,
	}

	return kitchenService
}

func (k *KitchenService) Run() {
	k.eventConsumer.RegisterEventHandler(events.RESTAURANT_EVENTS_TOPIC,
		events.RESTAURANT_CREATED_EVENT_ID,
		events.EventHandlerWrapper(k.restaurantEventHandler.HandleRestaurantCreatedEvent),
		reflect.TypeOf(events.RestaurantCreatedEvent{}))
	k.eventConsumer.RegisterEventHandler(events.RESTAURANT_EVENTS_TOPIC,
		events.RESTAURANT_DELETED_EVENT_ID,
		events.EventHandlerWrapper(k.restaurantEventHandler.HandleRestaurantDeletedEvent),
		reflect.TypeOf(events.RestaurantDeletedEvent{}))
	k.eventConsumer.RegisterEventHandler(events.RESTAURANT_EVENTS_TOPIC,
		events.MENU_ITEM_CREATED_EVENT_ID,
		events.EventHandlerWrapper(k.restaurantEventHandler.HandleMenuItemCreatedEvent),
		reflect.TypeOf(events.MenuItemCreatedEvent{}))
	k.eventConsumer.RegisterEventHandler(events.RESTAURANT_EVENTS_TOPIC,
		events.MENU_ITEM_UPDATED_EVENT_ID,
		events.EventHandlerWrapper(k.restaurantEventHandler.HandleMenuItemUpdatedEvent),
		reflect.TypeOf(events.MenuItemUpdatedEvent{}))
	k.eventConsumer.RegisterEventHandler(events.RESTAURANT_EVENTS_TOPIC,
		events.MENU_ITEM_DELETED_EVENT_ID,
		events.EventHandlerWrapper(k.restaurantEventHandler.HandleMenuItemDeletedEvent),
		reflect.TypeOf(events.MenuItemDeletedEvent{}))
	go k.eventConsumer.Run(context.Background())
}

func (k *KitchenService) Stop() {
	k.eventConsumerCancel()
	k.eventConsumer.Close()
}
