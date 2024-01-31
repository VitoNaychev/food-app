package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/order-svc/handlers"
	"github.com/VitoNaychev/food-app/order-svc/models"
	"github.com/VitoNaychev/food-app/pgconfig"
)

func main() {
	env := appenv.Enviornment{
		SecretKey: []byte(os.Getenv("SECRET")),

		Dbhost: "order-db",
		Dbport: "5432",
		Dbuser: os.Getenv("POSTGRES_USER"),
		Dbpass: os.Getenv("POSTGRES_PASSWORD"),
		Dbname: os.Getenv("POSTGRES_DB"),

		KafkaBrokers: strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
	}

	dbConfig := pgconfig.GetConfigFromEnv(env)
	connStr := dbConfig.GetConnectionString()

	orderStore, err := models.NewPgOrderStore(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Order Store error: %v", err)
	}

	orderItemStore, err := models.NewPgOrderItemStore(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Order Item Store error: %v", err)
	}

	addressStore, err := models.NewPgAddressStore(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Address Store error: %v", err)
	}

	eventPublisher, err := events.NewKafkaEventPublisher(env.KafkaBrokers)
	if err != nil {
		log.Fatalf("Kafka Event Publisher error: %v\n", err)
	}

	orderServer := handlers.NewOrderServer(orderStore, orderItemStore, addressStore, eventPublisher, handlers.VerifyJWT)

	fmt.Println("Order service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", orderServer))
}
