package handlers_test

import (
	"testing"
	"time"

	"github.com/VitoNaychev/food-app/delivery-svc/handlers"
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/delivery-svc/testdata"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/events/svcevents"
	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/VitoNaychev/food-app/testutil"
)

type StubDeliveryStore struct {
	deliveries      []models.Delivery
	updatedDelivery models.Delivery
}

func (d *StubDeliveryStore) GetDeliveryByID(id int) (models.Delivery, error) {
	for _, delivery := range d.deliveries {
		if delivery.ID == id {
			return delivery, nil
		}
	}

	return models.Delivery{}, storeerrors.ErrNotFound
}

func (d *StubDeliveryStore) UpdateDelivery(delivery *models.Delivery) error {
	d.updatedDelivery = *delivery

	return nil
}

func TestKitchenEventHandler(t *testing.T) {
	deliveryStore := &StubDeliveryStore{deliveries: []models.Delivery{testdata.VolenDelivery, testdata.PeterDelivery}}

	eventHandler := handlers.NewKitchenEventHandler(deliveryStore)

	t.Run("updates delivery state on TICKET_BEGIN_PREPARING event", func(t *testing.T) {
		readyBy, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 16:20:00")

		want := testdata.VolenDelivery
		want.State = models.IN_PROGRESS
		want.ReadyBy = readyBy

		// The AggregateID and the ID in the events are actually derived from the IDs in the
		// ticket's table in the kitchen svc, but because there is a 1:1 relation beetween
		// order IDs, ticket IDs and delivery IDs, the ID of the delivery is the same as the
		// ID of the associated ticket.
		payload := svcevents.TicketBeginPreparingEvent{
			ID:      want.ID,
			ReadyBy: readyBy,
		}
		event := events.NewTypedEvent(svcevents.TICKET_BEGIN_PREPARING_EVENT_ID, want.ID, payload)

		err := eventHandler.HandleTicketBeginPreparingEvent(event)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, deliveryStore.updatedDelivery, want)
	})

	t.Run("updates delivery state on TICKET_CANCEL event", func(t *testing.T) {
		want := testdata.VolenDelivery
		want.State = models.CANCELED

		payload := svcevents.TicketCancelEvent{
			ID: want.ID,
		}
		event := events.NewTypedEvent(svcevents.TICKET_CANCEL_EVENT_ID, want.ID, payload)

		err := eventHandler.HandleTicketCancelEvent(event)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, deliveryStore.updatedDelivery, want)
	})

	t.Run("updates delivery state on TICKET_FINISH_PREPARING event", func(t *testing.T) {
		want := testdata.PeterDelivery
		want.State = models.READY_FOR_PICKUP

		payload := svcevents.TicketFinishPreparingEvent{
			ID: want.ID,
		}
		event := events.NewTypedEvent(svcevents.TICKET_FINISH_PREPARING_EVENT_ID, want.ID, payload)

		err := eventHandler.HandleTicketFinishPreparingEvent(event)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, deliveryStore.updatedDelivery, want)
	})

}
