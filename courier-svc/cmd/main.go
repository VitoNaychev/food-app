package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/courier-svc/handlers"
	"github.com/VitoNaychev/food-app/courier-svc/models"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/pgconfig"
)

func main() {
	env := appenv.Enviornment{
		SecretKey: []byte(os.Getenv("SECRET")),
		ExpiresAt: 24 * time.Hour,

		Dbhost: "courier-db",
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
		log.Fatalf("Courier Store error: %v", err)
	}

	eventPublisher, err := events.NewKafkaEventPublisher(env.KafkaBrokers)
	if err != nil {
		log.Fatalf("Kafka Event Publisher error: %v\n", err)
	}

	server := handlers.NewCourierServer(env.SecretKey, env.ExpiresAt, &courierStore, eventPublisher)

	fmt.Println("courier service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", server))
}
