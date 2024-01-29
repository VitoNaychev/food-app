package services

import (
	"context"
	"testing"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/delivery-svc/handlers"
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/events"
)

type DeliveryService struct {
	CourierStore  *models.InMemoryCourierStore
	LocationStore *models.InMemoryLocationStore
	DeliveryStore *models.InMemoryDeliveryStore
	AddressStore  *models.InMemoryAddressStore

	CourierEventHandler *handlers.CourierEventHandler
	KitchenEventHandler *handlers.KitchenEventHandler

	EventConsumer       *events.KafkaEventConsumer
	EventConsumerCtx    context.Context
	EventConsumerCancel context.CancelFunc
}

func SetupDeliveryService(t testing.TB, env appenv.Enviornment, port string) DeliveryService {
	eventConsumer, err := events.NewKafkaEventConsumer(env.KafkaBrokers, "delivery-grp")
	if err != nil {
		t.Fatalf("Kafka Event Consumer error: %v\n", err)
	}

	courierStore := models.NewInMemoryCourierStore()
	locationStore := models.NewInMemoryLocationStore()
	deliveryStore := models.NewInMemoryDeliveryStore()
	addressStore := models.NewInMemoryAddressStore()

	courierEventHandler := handlers.NewCourierEventHandler(courierStore, locationStore)
	kitchenEventHandler := handlers.NewKitchenEventHandler(deliveryStore)
	eventConsumerCtx, eventConsumerCancel := context.WithCancel(context.Background())

	deliveryService := DeliveryService{
		CourierStore:  courierStore,
		LocationStore: locationStore,
		DeliveryStore: deliveryStore,
		AddressStore:  addressStore,

		CourierEventHandler: courierEventHandler,
		KitchenEventHandler: kitchenEventHandler,

		EventConsumer:       eventConsumer,
		EventConsumerCtx:    eventConsumerCtx,
		EventConsumerCancel: eventConsumerCancel,
	}

	return deliveryService
}

func (d *DeliveryService) Run() {
	handlers.RegisterCourierEventHandlers(d.EventConsumer, d.CourierEventHandler)
	handlers.RegisterKitchenEventHandlers(d.EventConsumer, d.KitchenEventHandler)

	go d.EventConsumer.Run(d.EventConsumerCtx)
	go events.LogEventConsumerErrors(d.EventConsumerCtx, d.EventConsumer)
}

func (d *DeliveryService) Stop() {
	d.EventConsumerCancel()
	d.EventConsumer.Close()
}
