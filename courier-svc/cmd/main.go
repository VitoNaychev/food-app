package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VitoNaychev/food-app/courier-svc/handlers"
	"github.com/VitoNaychev/food-app/courier-svc/models"
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
		postgresHost:     "courier-db",
		postgresPort:     "5432",
		postgresUser:     os.Getenv("POSTGRES_USER"),
		postgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		postgresDB:       os.Getenv("POSTGRES_DB"),
	}
	connStr := dbConfig.getConnectionString()

	courierStore, err := models.NewPgCourierStore(context.Background(), connStr)
	if err != nil {
		fmt.Printf("courier Store error: %v", err)
	}

	server := handlers.NewCourierServer(secretKey, expiresAt, &courierStore)

	fmt.Println("courier service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", server))
}
