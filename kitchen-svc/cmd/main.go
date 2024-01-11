package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/kitchen-svc/handlers"
	"github.com/VitoNaychev/food-app/kitchen-svc/models"
	"github.com/VitoNaychev/food-app/pgconfig"
)

func main() {
	env := appenv.Enviornment{
		SecretKey: []byte(os.Getenv("SECRET")),

		Dbhost: "kitchen-db",
		Dbport: "5432",
		Dbuser: os.Getenv("POSTGRES_USER"),
		Dbpass: os.Getenv("POSTGRES_PASSWORD"),
		Dbname: os.Getenv("POSTGRES_DB"),

		KafkaBrokers: strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
	}

	dbConfig := pgconfig.GetConfigFromEnv(env)
	connStr := dbConfig.GetConnectionString()

	restaurantStore, err := models.NewPgRestaurantStore(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Restaurant Store error: %v\n", err)
	}

	menuItemStore, err := models.NewPgMenuItemStore(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Menu Item Store error: %v\n", err)
	}

	ticketStore, err := models.NewPgTicketStore(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Ticket Store error: %v\n", err)
	}

	ticketItemStore, err := models.NewPgTicketItemStore(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Ticket Item Store error: %v\n", err)
	}

	eventConsumer, err := events.NewKafkaEventConsumer(env.KafkaBrokers, "kitchen-svc")
	if err != nil {
		log.Fatalf("Kafka Event Consumer error: %v\n", err)
	}

	restaurantEventhandler := handlers.NewRestaurantEventHandler(restaurantStore, menuItemStore)

	registerRestaurantEventHandlers(eventConsumer, *restaurantEventhandler)
	go eventConsumer.Run(context.Background())
	go events.LogEventConsumerErrors(context.Background(), eventConsumer)

	kitchenServer := handlers.NewTicketServer(env.SecretKey, ticketStore, ticketItemStore, menuItemStore, restaurantStore)

	log.Println("kitchen service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", kitchenServer))
}

func registerRestaurantEventHandlers(eventConsumer events.EventConsumer, restaurantEventHandler handlers.RestaurantEventHandler) {
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
