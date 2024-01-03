package svcintegration

import (
	"testing"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/kitchen-svc/handlers"
	"github.com/VitoNaychev/food-app/kitchen-svc/models"
)

type KitchenService struct {
	restaurantStore models.RestaurantStore

	restaurantEventHandler *handlers.RestaurantEventHandler

	eventConsumer *events.KafkaEventConsumer
}

func SetupKitchenService(t testing.TB, env appenv.Enviornment, port string) KitchenService {
	eventConsumer, err := events.NewKafkaEventConsumer(env.KafkaBrokers, "kitchen-grp")
	if err != nil {
		t.Fatalf("Kafka Event Consumer error: %v\n", err)
	}

	restaurantStore := models.NewInMemoryRestaurantStore()
	restaurantEventHandler := handlers.NewRestaurantEventHandler(restaurantStore)

	kitchenService := KitchenService{
		restaurantStore: restaurantStore,

		restaurantEventHandler: restaurantEventHandler,

		eventConsumer: eventConsumer,
	}

	return kitchenService
}

func (k *KitchenService) Run() {
	k.eventConsumer.RegisterEventHandler(events.RESTAURANT_EVENTS_TOPIC, k.restaurantEventHandler.HandleRestaurantEvent)
}

func (k *KitchenService) Stop() {
	k.eventConsumer.Close()
}
