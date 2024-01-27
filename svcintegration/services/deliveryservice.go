package services

import (
	"context"
	"reflect"
	"testing"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/delivery-svc/handlers"
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/events/svcevents"
)

type DeliveryService struct {
	CourierStore  *models.InMemoryCourierStore
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
	deliveryStore := models.NewInMemoryDeliveryStore()
	addressStore := models.NewInMemoryAddressStore()

	courierEventHandler := handlers.NewCourierEventHandler(courierStore)
	kitchenEventHandler := handlers.NewKitchenEventHandler(deliveryStore)
	eventConsumerCtx, eventConsumerCancel := context.WithCancel(context.Background())

	deliveryService := DeliveryService{
		CourierStore:  courierStore,
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
	// d.EventConsumer.RegisterEventHandler(svcevents.COURIER_EVENTS_TOPIC,
	// 	svcevents.COURIER_CREATED_EVENT_ID,
	// 	events.EventHandlerWrapper(d.CourierEventHandler.HandleCourierCreatedEvent),
	// 	reflect.TypeOf(svcevents.CourierCreatedEvent{}))
	// d.EventConsumer.RegisterEventHandler(svcevents.COURIER_EVENTS_TOPIC,
	// 	svcevents.COURIER_DELETED_EVENT_ID,
	// 	events.EventHandlerWrapper(d.CourierEventHandler.HandleCourierDeletedEvent),
	// 	reflect.TypeOf(svcevents.CourierDeletedEvent{}))

	d.EventConsumer.RegisterEventHandler(svcevents.KITCHEN_EVENTS_TOPIC,
		svcevents.TICKET_BEGIN_PREPARING_EVENT_ID,
		events.EventHandlerWrapper(d.KitchenEventHandler.HandleTicketBeginPreparingEvent),
		reflect.TypeOf(svcevents.TicketBeginPreparingEvent{}))
	d.EventConsumer.RegisterEventHandler(svcevents.KITCHEN_EVENTS_TOPIC,
		svcevents.TICKET_FINISH_PREPARING_EVENT_ID,
		events.EventHandlerWrapper(d.KitchenEventHandler.HandleTicketFinishPreparingEvent),
		reflect.TypeOf(svcevents.TicketFinishPreparingEvent{}))

	go d.EventConsumer.Run(d.EventConsumerCtx)
	go events.LogEventConsumerErrors(d.EventConsumerCtx, d.EventConsumer)
}

func (d *DeliveryService) Stop() {
	d.EventConsumerCancel()
	d.EventConsumer.Close()
}
