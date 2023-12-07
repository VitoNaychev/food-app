package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VitoNaychev/food-app/restaurant-svc/handlers"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
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
		postgresHost:     "restaurant-db",
		postgresPort:     "5432",
		postgresUser:     os.Getenv("POSTGRES_USER"),
		postgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		postgresDB:       os.Getenv("POSTGRES_DB"),
	}
	connStr := dbConfig.getConnectionString()

	restaurantStore, err := models.NewPgRestaurantStore(context.Background(), connStr)
	if err != nil {
		fmt.Printf("Restaurant Store error: %v", err)
	}

	addressStore, err := models.NewPgAddressStore(context.Background(), connStr)
	if err != nil {
		fmt.Printf("Address Store error: %v", err)
	}

	hoursStore, err := models.NewPgHoursStore(context.Background(), connStr)
	if err != nil {
		fmt.Printf("Working Hours Store error: %v", err)
	}

	menuStore, err := models.NewPgMenuStore(context.Background(), connStr)
	if err != nil {
		fmt.Printf("Menu Store error: %v", err)
	}

	restaurantServer := handlers.NewRestaurantServer(secretKey, expiresAt, &restaurantStore)
	addressServer := handlers.NewAddressServer(secretKey, &addressStore, &restaurantStore)
	hoursServer := handlers.NewHoursServer(secretKey, &hoursStore, &restaurantStore)
	menuServer := handlers.NewMenuServer(secretKey, &menuStore, &restaurantStore)

	router := handlers.NewRouterServer(restaurantServer, addressServer, hoursServer, menuServer)

	fmt.Println("Restaurant service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
