package handlers_test

import (
	"testing"

	"github.com/VitoNaychev/food-app/delivery-svc/handlers"
	"github.com/VitoNaychev/food-app/delivery-svc/stubs"
	"github.com/VitoNaychev/food-app/delivery-svc/testdata"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/events/svcevents"
	"github.com/VitoNaychev/food-app/testutil"
)

func TestCourierEventHandler(t *testing.T) {
	courierStore := &stubs.StubCourierStore{}

	eventHandler := handlers.NewCourierEventHandler(courierStore)

	t.Run("creates courier on COURIER_CREATED_EVENT", func(t *testing.T) {
		want := testdata.VolenCourier

		payload := svcevents.CourierCreatedEvent{
			ID:   want.ID,
			Name: want.Name,
		}
		event := events.NewTypedEvent(svcevents.COURIER_CREATED_EVENT_ID, 1, payload)

		err := eventHandler.HandleCourierCreatedEvent(event)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, courierStore.CreatedCourier, want)
	})

	t.Run("deletes courier on COURIER_DELETED_EVENT", func(t *testing.T) {
		want := testdata.VolenCourier

		payload := svcevents.CourierDeletedEvent{
			ID: want.ID,
		}
		event := events.NewTypedEvent(svcevents.COURIER_DELETED_EVENT_ID, 1, payload)

		err := eventHandler.HandleCourierDeletedEvent(event)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, courierStore.DeletedCourierID, want.ID)
	})
}
