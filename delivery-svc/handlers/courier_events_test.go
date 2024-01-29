package handlers_test

import (
	"testing"

	"github.com/VitoNaychev/food-app/delivery-svc/handlers"
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/delivery-svc/stubs"
	"github.com/VitoNaychev/food-app/delivery-svc/testdata"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/events/svcevents"
	"github.com/VitoNaychev/food-app/testutil"
)

func TestCourierEventHandler(t *testing.T) {
	courierStore := &stubs.StubCourierStore{}
	locationStore := &stubs.StubLocationStore{}

	eventHandler := handlers.NewCourierEventHandler(courierStore, locationStore)

	t.Run("creates courier and location on COURIER_CREATED_EVENT", func(t *testing.T) {
		wantCourier := testdata.VolenCourier
		wantLocation := models.Location{
			CourierID: wantCourier.ID,
			Lat:       0.0,
			Lon:       0.0,
		}

		payload := svcevents.CourierCreatedEvent{
			ID:   wantCourier.ID,
			Name: wantCourier.Name,
		}
		event := events.NewTypedEvent(svcevents.COURIER_CREATED_EVENT_ID, 1, payload)

		err := eventHandler.HandleCourierCreatedEvent(event)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, courierStore.CreatedCourier, wantCourier)
		testutil.AssertEqual(t, locationStore.CreatedLocation, wantLocation)
	})

	t.Run("deletes courier and location on COURIER_DELETED_EVENT", func(t *testing.T) {
		want := testdata.VolenCourier

		payload := svcevents.CourierDeletedEvent{
			ID: want.ID,
		}
		event := events.NewTypedEvent(svcevents.COURIER_DELETED_EVENT_ID, 1, payload)

		err := eventHandler.HandleCourierDeletedEvent(event)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, courierStore.DeletedCourierID, want.ID)
		testutil.AssertEqual(t, locationStore.DeletedLocationCourierID, want.ID)
	})
}
