package handlers

import (
	"reflect"

	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/events/svcevents"
	"github.com/VitoNaychev/food-app/kitchen-svc/models"
)

type OrderEventHandler struct {
	ticketStore     models.TicketStore
	ticketItemStore models.TicketItemStore
}

func NewOrderEventHandler(ticketStore models.TicketStore, ticketItemStore models.TicketItemStore) *OrderEventHandler {
	orderEventHandler := OrderEventHandler{
		ticketStore:     ticketStore,
		ticketItemStore: ticketItemStore,
	}

	return &orderEventHandler
}

func RegisterOrderEventHandlers(eventConsumer events.EventConsumer, orderEventHandler *OrderEventHandler) {
	eventConsumer.RegisterEventHandler(svcevents.ORDER_EVENTS_TOPIC,
		svcevents.ORDER_CREATED_EVENT_ID,
		events.EventHandlerWrapper(orderEventHandler.HandleOrderCreatedEvent),
		reflect.TypeOf(svcevents.OrderCreatedEvent{}))
}

func (o *OrderEventHandler) HandleOrderCreatedEvent(event events.Event[svcevents.OrderCreatedEvent]) error {
	ticket := TicketFromOrderCreatedEvent(event.Payload)
	err := o.ticketStore.CreateTicket(&ticket)
	if err != nil {
		return err
	}

	ticketItems := TicketItemsFromOrderCreatedEvent(event.Payload)
	for _, ticketItem := range ticketItems {
		err := o.ticketItemStore.CreateTicketItem(&ticketItem)
		if err != nil {
			return err
		}
	}

	return nil
}

func TicketFromOrderCreatedEvent(orderCreatedEvent svcevents.OrderCreatedEvent) models.Ticket {
	ticket := models.Ticket{
		ID:           orderCreatedEvent.ID,
		State:        models.OPEN,
		RestaurantID: orderCreatedEvent.RestaurantID,
		Total:        orderCreatedEvent.Total,
		ReadyBy:      models.ZeroTime,
	}

	return ticket
}

func TicketItemsFromOrderCreatedEvent(orderCreatedEvent svcevents.OrderCreatedEvent) []models.TicketItem {
	ticketItems := []models.TicketItem{}

	for _, orderCreatedEventItem := range orderCreatedEvent.Items {
		ticketItem := models.TicketItem{
			ID:         orderCreatedEventItem.ID,
			TicketID:   orderCreatedEvent.ID,
			MenuItemID: orderCreatedEventItem.MenuItemID,
			Quantity:   orderCreatedEventItem.Quantity,
		}
		ticketItems = append(ticketItems, ticketItem)
	}

	return ticketItems
}
