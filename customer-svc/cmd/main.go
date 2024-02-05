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
	"github.com/VitoNaychev/food-app/customer-svc/handlers"
	"github.com/VitoNaychev/food-app/customer-svc/models"
	"github.com/VitoNaychev/food-app/pgconfig"
)

func main() {
	env := appenv.Enviornment{
		SecretKey: []byte(os.Getenv("SECRET")),
		ExpiresAt: 24 * time.Hour,

		Dbhost: "customer-db",
		Dbport: "5432",
		Dbuser: os.Getenv("POSTGRES_USER"),
		Dbpass: os.Getenv("POSTGRES_PASSWORD"),
		Dbname: os.Getenv("POSTGRES_DB"),

		KafkaBrokers: strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
	}

	dbConfig := pgconfig.GetConfigFromEnv(env)
	connStr := dbConfig.GetConnectionString()

	customerStore, err := models.NewPgCustomerStore(context.Background(), connStr)
	if err != nil {
		fmt.Printf("Customer Store error: %v", err)
	}

	addressStore, err := models.NewPgAddressStore(context.Background(), connStr)
	if err != nil {
		fmt.Printf("Address Store error: %v", err)
	}

	customerServer := handlers.NewCustomerServer(env.SecretKey, env.ExpiresAt, &customerStore)
	addressServer := handlers.NewCustomerAddressServer(&addressStore, &customerStore, env.SecretKey)

	router := handlers.NewRouterServer(customerServer, addressServer)

	fmt.Println("Customer service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
