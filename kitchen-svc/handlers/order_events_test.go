package handlers_test

import (
	"testing"

	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/events/svcevents"
	"github.com/VitoNaychev/food-app/kitchen-svc/handlers"
	"github.com/VitoNaychev/food-app/kitchen-svc/models"
	"github.com/VitoNaychev/food-app/kitchen-svc/stubs"
	"github.com/VitoNaychev/food-app/kitchen-svc/testdata"

	"github.com/VitoNaychev/food-app/testutil"
)

func TestOrderCreatedEventHandler(t *testing.T) {
	ticketStore := &stubs.StubTicketStore{}
	ticketItemStore := &stubs.StubTicketItemStore{}

	eventHandler := handlers.NewOrderEventHandler(ticketStore, ticketItemStore)

	t.Run("creates corresponding delivery", func(t *testing.T) {
		wantTicket := testdata.OpenShackTicket
		wantTicketItems := []models.TicketItem{testdata.OpenShackTicketItems}

		payload := testdata.PeterOrderCreatedEvent
		event := events.NewTypedEvent(svcevents.ORDER_CREATED_EVENT_ID, testdata.PeterOrderCreatedEvent.ID, payload)

		err := eventHandler.HandleOrderCreatedEvent(event)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, ticketStore.SpyTicket, wantTicket)
		testutil.AssertEqual(t, ticketItemStore.SpyTicketItems, wantTicketItems)
	})
}
