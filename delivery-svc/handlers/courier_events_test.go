package handlers_test

import (
	"testing"

	"github.com/VitoNaychev/food-app/delivery-svc/handlers"
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/delivery-svc/testdata"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/events/svcevents"
	"github.com/VitoNaychev/food-app/testutil"
)

type StubCourierStore struct {
	createdCourier   models.Courier
	deletedCourierID int
}

func (s *StubCourierStore) CreateCourier(courier *models.Courier) error {
	s.createdCourier = *courier
	return nil
}

func (s *StubCourierStore) DeleteCourier(id int) error {
	s.deletedCourierID = id
	return nil
}

func TestCourierEventHandler(t *testing.T) {
	courierStore := &StubCourierStore{}

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
		testutil.AssertEqual(t, courierStore.createdCourier, want)
	})

	t.Run("deletes courier on COURIER_DELETED_EVENT", func(t *testing.T) {
		want := testdata.VolenCourier

		payload := svcevents.CourierDeletedEvent{
			ID: want.ID,
		}
		event := events.NewTypedEvent(svcevents.COURIER_DELETED_EVENT_ID, 1, payload)

		err := eventHandler.HandleCourierDeletedEvent(event)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, courierStore.deletedCourierID, want.ID)
	})
}
