package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/delivery-svc/handlers"
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/pgconfig"
)

func main() {
	env := appenv.Enviornment{
		SecretKey: []byte(os.Getenv("SECRET")),

		Dbhost: "delivery-db",
		Dbport: "5432",
		Dbuser: os.Getenv("POSTGRES_USER"),
		Dbpass: os.Getenv("POSTGRES_PASSWORD"),
		Dbname: os.Getenv("POSTGRES_DB"),

		KafkaBrokers: strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
	}

	dbConfig := pgconfig.GetConfigFromEnv(env)
	connStr := dbConfig.GetConnectionString()

	courierStore, err := models.NewPgCourierStore(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Courier Store error: %v\n", err)
	}

	locationStore, err := models.NewPgLocationStore(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Location Store error: %v\n", err)
	}

	addressStore, err := models.NewPgAddressStore(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Address Store error: %v\n", err)
	}

	deliveryStore, err := models.NewPgDeliveryStore(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Delivery Store error: %v\n", err)
	}

	eventConsumer, err := events.NewKafkaEventConsumer(env.KafkaBrokers, "delivery-svc")
	if err != nil {
		log.Fatalf("Kafka Event Consumer error: %v\n", err)
	}

	courierEventHandler := handlers.NewCourierEventHandler(courierStore, locationStore)
	handlers.RegisterCourierEventHandlers(eventConsumer, courierEventHandler)

	kitchenEventHandler := handlers.NewKitchenEventHandler(deliveryStore)
	handlers.RegisterKitchenEventHandlers(eventConsumer, kitchenEventHandler)

	orderEventHandler := handlers.NewOrderEventHandler(deliveryStore, addressStore)
	handlers.RegisterOrderEventHandlers(eventConsumer, orderEventHandler)

	go eventConsumer.Run(context.Background())
	go events.LogEventConsumerErrors(context.Background(), eventConsumer)

	locationServer := handlers.NewLocationServer(env.SecretKey, locationStore, courierStore)
	deliveryServer := handlers.NewDeliveryServer(env.SecretKey, deliveryStore, addressStore, courierStore)

	router := handlers.NewRouterServer(deliveryServer, locationServer)

	log.Println("Delivery service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
