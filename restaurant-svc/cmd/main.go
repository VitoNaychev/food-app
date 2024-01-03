package main

import (
	"context"
	"os"
	"time"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/restaurant-svc/service"
)

func main() {
	env := appenv.Enviornment{
		SecretKey: []byte(os.Getenv("SECRET")),
		ExpiresAt: 24 * time.Hour,

		Dbhost: "restaurant-db",
		Dbport: "5432",
		Dbuser: os.Getenv("POSTGRES_USER"),
		Dbpass: os.Getenv("POSTGRES_PASSWORD"),
		Dbname: os.Getenv("POSTGRES_DB"),
	}

	service.Run(context.Background(), env, ":8080")
}
