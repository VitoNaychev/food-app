package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/pgconfig"
	"github.com/VitoNaychev/food-app/restaurant-svc/handlers"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

func Run(ctx context.Context, env appenv.Enviornment, port string) {
	dbConfig := pgconfig.GetConfigFromEnv(env)
	connStr := dbConfig.GetConnectionString()

	restaurantStore, err := models.NewPgRestaurantStore(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Restaurant Store error: %v\n", err)
	}

	addressStore, err := models.NewPgAddressStore(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Address Store error: %v\n", err)
	}

	hoursStore, err := models.NewPgHoursStore(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Working Hours Store error: %v\n", err)
	}

	menuStore, err := models.NewPgMenuStore(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Menu Store error: %v\n", err)
	}

	eventPublisher, err := events.NewKafkaEventPublisher(env.KafkaBrokers)
	if err != nil {
		log.Fatalf("Event Publisher error: %v\n", err)
	}
	defer eventPublisher.Close()

	restaurantServer := handlers.NewRestaurantServer(env.SecretKey, env.ExpiresAt, &restaurantStore, eventPublisher)
	addressServer := handlers.NewAddressServer(env.SecretKey, &addressStore, &restaurantStore)
	hoursServer := handlers.NewHoursServer(env.SecretKey, &hoursStore, &restaurantStore)
	menuServer := handlers.NewMenuServer(env.SecretKey, &menuStore, &restaurantStore, eventPublisher)

	router := handlers.NewRouterServer(restaurantServer, addressServer, hoursServer, menuServer)

	server := http.Server{
		Addr:    port,
		Handler: router,
	}

	fmt.Printf("Restaurant service listening on %s\n", port)
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v\n", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, time.Second)
	defer shutdownCancel()

	err = server.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatalf("Shutdown error: %v\n", err)
	}
}
