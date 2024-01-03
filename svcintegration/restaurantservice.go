package svcintegration

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/restaurant-svc/handlers"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

type RestaurantService struct {
	restaurantStore models.RestaurantStore
	addressStore    models.AddressStore
	hoursStore      models.HoursStore
	menuStore       models.MenuStore

	restaurantHandler *handlers.RestaurantServer
	addressHandler    *handlers.AddressServer
	hoursHandler      *handlers.HoursServer
	menuHandler       *handlers.MenuServer

	router *handlers.RouterServer

	server *http.Server

	eventPublisher *events.KafkaEventPublisher
}

func SetupRestaurantService(t testing.TB, env appenv.Enviornment, port string) RestaurantService {
	eventPublisher, err := events.NewKafkaEventPublisher(env.KafkaBrokers)
	if err != nil {
		t.Fatalf("Event Publisher error: %v\n", err)
	}

	restaurantStore := models.NewInMemoryRestaurantStore()
	addressStore := models.NewInMemoryAddressStore()
	hoursStore := models.NewInMemoryHoursStore()
	menuStore := models.NewInMemoryMenuStore()

	restaurantHandler := handlers.NewRestaurantServer(env.SecretKey, env.ExpiresAt, restaurantStore, eventPublisher)
	addressHandler := handlers.NewAddressServer(env.SecretKey, addressStore, restaurantStore)
	hoursHandler := handlers.NewHoursServer(env.SecretKey, hoursStore, restaurantStore)
	menuHandler := handlers.NewMenuServer(env.SecretKey, menuStore, restaurantStore, eventPublisher)

	router := handlers.NewRouterServer(restaurantHandler, addressHandler, hoursHandler, menuHandler)

	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	restaurantService := RestaurantService{
		restaurantStore: restaurantStore,
		addressStore:    addressStore,
		hoursStore:      hoursStore,
		menuStore:       menuStore,

		restaurantHandler: restaurantHandler,
		addressHandler:    addressHandler,
		hoursHandler:      hoursHandler,
		menuHandler:       menuHandler,

		router: router,

		server: server,

		eventPublisher: eventPublisher,
	}

	return restaurantService
}

func (r *RestaurantService) Run() {
	log.Printf("Restaurant service listening on %s\n", r.server.Addr)

	go func() {
		err := r.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v\n", err)
		}
	}()
}

func (r *RestaurantService) Stop() {
	r.eventPublisher.Close()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second)
	defer shutdownCancel()

	err := r.server.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatalf("Shutdown error: %v\n", err)
	}
}
