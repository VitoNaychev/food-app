package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/VitoNaychev/bt-order-svc/handlers"
	"github.com/VitoNaychev/bt-order-svc/models"
)

type DBConfig struct {
	postgresHost     string
	postgresPort     string
	postgresUser     string
	postgresPassword string
	postgresDB       string
}

func (d *DBConfig) getConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		d.postgresUser, d.postgresPassword, d.postgresHost, d.postgresPort, d.postgresDB)
}

func main() {
	dbConfig := DBConfig{
		postgresHost:     "order-db",
		postgresPort:     "5432",
		postgresUser:     os.Getenv("POSTGRES_USER"),
		postgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		postgresDB:       os.Getenv("POSTGRES_DB"),
	}
	connStr := dbConfig.getConnectionString()

	orderStore, err := models.NewPgOrderStore(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Order Store error: %v", err)
	}

	addressStore, err := models.NewPgAddressStore(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Address Store error: %v", err)
	}

	orderServer := handlers.NewOrderServer(orderStore, addressStore, handlers.VerifyJWT)

	fmt.Println("Order service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", orderServer))
}
