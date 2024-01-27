package services

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/courier-svc/handlers"
	"github.com/VitoNaychev/food-app/courier-svc/models"
	"github.com/VitoNaychev/food-app/events"
)

type CourierService struct {
	CourierStore models.CourierStore

	CourierHandler *handlers.CourierServer

	Server *http.Server

	EventPublisher *events.KafkaEventPublisher
}

func SetupCourierService(t testing.TB, env appenv.Enviornment, port string) CourierService {
	eventPublisher, err := events.NewKafkaEventPublisher(env.KafkaBrokers)
	if err != nil {
		t.Fatalf("Event Publisher error: %v\n", err)
	}

	courierStore := models.NewInMemoryCourierStore()

	courierHandler := handlers.NewCourierServer(env.SecretKey, env.ExpiresAt, courierStore, eventPublisher)

	server := &http.Server{
		Addr:    port,
		Handler: courierHandler,
	}

	courierService := CourierService{
		CourierStore: courierStore,

		CourierHandler: courierHandler,

		Server: server,

		EventPublisher: eventPublisher,
	}

	return courierService
}

func (c *CourierService) Run() {
	log.Printf("Courier service listening on %s\n", c.Server.Addr)

	go func() {
		err := c.Server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v\n", err)
		}
	}()
}

func (c *CourierService) Stop() {
	c.EventPublisher.Close()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second)
	defer shutdownCancel()

	err := c.Server.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatalf("Shutdown error: %v\n", err)
	}
}
