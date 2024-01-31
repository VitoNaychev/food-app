package main

import (
	"context"
	"log"
	"net/http"
	"os"
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

	eventPublisher, err := events.NewKafkaEventPublisher(env.KafkaBrokers)
	if err != nil {
		log.Fatalf("Kafka Event Publisher error: %v\n", err)
	}

	eventConsumer, err := events.NewKafkaEventConsumer(env.KafkaBrokers, "kitchen-svc")
	if err != nil {
		log.Fatalf("Kafka Event Consumer error: %v\n", err)
	}

	restaurantEventhandler := handlers.NewRestaurantEventHandler(restaurantStore, menuItemStore)
	orderEventHandler := handlers.NewOrderEventHandler(ticketStore, ticketItemStore)

	handlers.RegisterRestaurantEventHandlers(eventConsumer, restaurantEventhandler)
	handlers.RegisterOrderEventHandlers(eventConsumer, orderEventHandler)
	go eventConsumer.Run(context.Background())
	go events.LogEventConsumerErrors(context.Background(), eventConsumer)

	kitchenServer := handlers.NewTicketServer(env.SecretKey, ticketStore, ticketItemStore, menuItemStore, restaurantStore, eventPublisher)

	log.Println("kitchen service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", kitchenServer))
}
