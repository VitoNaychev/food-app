package services

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
	RestaurantStore models.RestaurantStore
	AddressStore    models.AddressStore
	HoursStore      models.HoursStore
	MenuStore       models.MenuStore

	RestaurantHandler *handlers.RestaurantServer
	AddressHandler    *handlers.AddressServer
	HoursHandler      *handlers.HoursServer
	MenuHandler       *handlers.MenuServer

	Router *handlers.RouterServer

	Server *http.Server

	EventPublisher *events.KafkaEventPublisher
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
		RestaurantStore: restaurantStore,
		AddressStore:    addressStore,
		HoursStore:      hoursStore,
		MenuStore:       menuStore,

		RestaurantHandler: restaurantHandler,
		AddressHandler:    addressHandler,
		HoursHandler:      hoursHandler,
		MenuHandler:       menuHandler,

		Router: router,

		Server: server,

		EventPublisher: eventPublisher,
	}

	return restaurantService
}

func (r *RestaurantService) Run() {
	log.Printf("Restaurant service listening on %s\n", r.Server.Addr)

	go func() {
		err := r.Server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v\n", err)
		}
	}()
}

func (r *RestaurantService) Stop() {
	r.EventPublisher.Close()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second)
	defer shutdownCancel()

	err := r.Server.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatalf("Shutdown error: %v\n", err)
	}
}
