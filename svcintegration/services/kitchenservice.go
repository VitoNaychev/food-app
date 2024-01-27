package services

import (
	"context"
	"log"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/kitchen-svc/handlers"
	"github.com/VitoNaychev/food-app/kitchen-svc/models"
)

type KitchenService struct {
	RestaurantStore models.RestaurantStore
	MenuItemStore   models.MenuItemStore
	TicketStore     *models.InMemoryTicketStore
	TicketItemStore *models.InMemoryTicketItemStore

	Server *http.Server

	RestaurantEventHandler *handlers.RestaurantEventHandler

	EventPublisher *events.KafkaEventPublisher

	EventConsumer       *events.KafkaEventConsumer
	EventConsumerCtx    context.Context
	EventConsumerCancel context.CancelFunc
}

func SetupKitchenService(t testing.TB, env appenv.Enviornment, port string) KitchenService {
	eventPublisher, err := events.NewKafkaEventPublisher(env.KafkaBrokers)
	if err != nil {
		t.Fatalf("Kafka Event Publisher error: %v\n", err)
	}

	eventConsumer, err := events.NewKafkaEventConsumer(env.KafkaBrokers, "kitchen-grp")
	if err != nil {
		t.Fatalf("Kafka Event Consumer error: %v\n", err)
	}

	restaurantStore := models.NewInMemoryRestaurantStore()
	menuItemStore := models.NewInMemoryMenuItemStore()
	ticketStore := models.NewInMemoryTicketStore()
	ticketItemStore := models.NewInMemoryTicketItemStore()

	ticketHandler := handlers.NewTicketServer(env.SecretKey, ticketStore, ticketItemStore, menuItemStore, restaurantStore, eventPublisher)
	server := &http.Server{
		Addr:    port,
		Handler: ticketHandler,
	}

	restaurantEventHandler := handlers.NewRestaurantEventHandler(restaurantStore, menuItemStore)
	eventConsumerCtx, eventConsumerCancel := context.WithCancel(context.Background())

	kitchenService := KitchenService{
		RestaurantStore: restaurantStore,

		RestaurantEventHandler: restaurantEventHandler,
		MenuItemStore:          menuItemStore,
		TicketStore:            ticketStore,
		TicketItemStore:        ticketItemStore,

		Server: server,

		EventPublisher: eventPublisher,

		EventConsumer:       eventConsumer,
		EventConsumerCtx:    eventConsumerCtx,
		EventConsumerCancel: eventConsumerCancel,
	}

	return kitchenService
}

func (k *KitchenService) Run() {
	k.EventConsumer.RegisterEventHandler(events.RESTAURANT_EVENTS_TOPIC,
		events.RESTAURANT_CREATED_EVENT_ID,
		events.EventHandlerWrapper(k.RestaurantEventHandler.HandleRestaurantCreatedEvent),
		reflect.TypeOf(events.RestaurantCreatedEvent{}))
	k.EventConsumer.RegisterEventHandler(events.RESTAURANT_EVENTS_TOPIC,
		events.RESTAURANT_DELETED_EVENT_ID,
		events.EventHandlerWrapper(k.RestaurantEventHandler.HandleRestaurantDeletedEvent),
		reflect.TypeOf(events.RestaurantDeletedEvent{}))
	k.EventConsumer.RegisterEventHandler(events.RESTAURANT_EVENTS_TOPIC,
		events.MENU_ITEM_CREATED_EVENT_ID,
		events.EventHandlerWrapper(k.RestaurantEventHandler.HandleMenuItemCreatedEvent),
		reflect.TypeOf(events.MenuItemCreatedEvent{}))
	k.EventConsumer.RegisterEventHandler(events.RESTAURANT_EVENTS_TOPIC,
		events.MENU_ITEM_UPDATED_EVENT_ID,
		events.EventHandlerWrapper(k.RestaurantEventHandler.HandleMenuItemUpdatedEvent),
		reflect.TypeOf(events.MenuItemUpdatedEvent{}))
	k.EventConsumer.RegisterEventHandler(events.RESTAURANT_EVENTS_TOPIC,
		events.MENU_ITEM_DELETED_EVENT_ID,
		events.EventHandlerWrapper(k.RestaurantEventHandler.HandleMenuItemDeletedEvent),
		reflect.TypeOf(events.MenuItemDeletedEvent{}))
	go k.EventConsumer.Run(k.EventConsumerCtx)

	log.Printf("Kitchen service listening on %s\n", k.Server.Addr)

	go func() {
		err := k.Server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v\n", err)
		}
	}()
}

func (k *KitchenService) Stop() {
	k.EventPublisher.Close()

	k.EventConsumerCancel()
	k.EventConsumer.Close()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second)
	defer shutdownCancel()

	err := k.Server.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatalf("Shutdown error: %v\n", err)
	}
}
