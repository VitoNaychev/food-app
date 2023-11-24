package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VitoNaychev/bt-customer-svc/handlers"
	"github.com/VitoNaychev/bt-customer-svc/models"
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
	secretKey := []byte(os.Getenv("SECRET"))
	expiresAt := 24 * time.Hour

	dbConfig := DBConfig{
		postgresHost:     "customer-db",
		postgresPort:     "5432",
		postgresUser:     os.Getenv("POSTGRES_USER"),
		postgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		postgresDB:       os.Getenv("POSTGRES_DB"),
	}
	connStr := dbConfig.getConnectionString()

	customerStore, err := models.NewPgCustomerStore(context.Background(), connStr)
	if err != nil {
		fmt.Printf("Customer Store error: %v", err)
	}

	addressStore, err := models.NewPgAddressStore(context.Background(), connStr)
	if err != nil {
		fmt.Printf("Address Store error: %v", err)
	}

	customerServer := handlers.NewCustomerServer(secretKey, expiresAt, &customerStore)
	addressServer := handlers.NewCustomerAddressServer(&addressStore, &customerStore, secretKey)

	router := handlers.NewRouterServer(customerServer, addressServer)

	fmt.Println("Customer service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
