package services

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/msgtypes"
	"github.com/VitoNaychev/food-app/order-svc/handlers"
	"github.com/VitoNaychev/food-app/order-svc/models"
)

func dummyVerifyJWT(jwt string) (msgtypes.AuthResponse, error) {
	id, _ := strconv.Atoi(jwt)
	return msgtypes.AuthResponse{Status: msgtypes.OK, ID: id}, nil
}

type OrderService struct {
	OrderStore     *models.InMemoryOrderStore
	OrderItemStore *models.InMemoryOrderItemStore
	AddressStore   *models.InMemoryAddressStore

	OrderHandler handlers.OrderServer

	Server *http.Server

	EventPublisher *events.KafkaEventPublisher
}

func SetupOrderService(t testing.TB, env appenv.Enviornment, port string) OrderService {
	eventPublisher, err := events.NewKafkaEventPublisher(env.KafkaBrokers)
	if err != nil {
		t.Fatalf("Kafka Event Publisher error: %v\n", err)
	}

	orderStore := models.NewInMemoryOrderStore()
	orderItemStore := models.NewInMemoryOrderItemStore()
	addressStore := models.NewInMemoryAddressStore()

	orderHandler := handlers.NewOrderServer(orderStore, orderItemStore, addressStore, eventPublisher, dummyVerifyJWT)

	server := &http.Server{
		Addr:    port,
		Handler: orderHandler,
	}

	orderService := OrderService{
		OrderStore:     orderStore,
		OrderItemStore: orderItemStore,
		AddressStore:   addressStore,

		OrderHandler: orderHandler,

		Server: server,

		EventPublisher: eventPublisher,
	}

	return orderService
}

func (o *OrderService) Run() {
	log.Printf("Order service listening on %s\n", o.Server.Addr)

	go func() {
		err := o.Server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v\n", err)
		}
	}()
}

func (o *OrderService) Stop() {
	o.EventPublisher.Close()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second)
	defer shutdownCancel()

	err := o.Server.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatalf("Shutdown error: %v\n", err)
	}
}
