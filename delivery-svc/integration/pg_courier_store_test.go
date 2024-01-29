package integration

import (
	"context"
	"testing"

	"github.com/VitoNaychev/food-app/delivery-svc/handlers"
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/delivery-svc/testdata"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/events/svcevents"
	"github.com/VitoNaychev/food-app/integrationutil"
	"github.com/VitoNaychev/food-app/pgconfig"
	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/VitoNaychev/food-app/testutil"
)

func TestCourierEventHandlerIntegration(t *testing.T) {
	config := pgconfig.GetConfigFromEnv(env)
	integrationutil.SetupDatabaseContainer(t, &config, "../sql-scripts/init.sql")

	connStr := config.GetConnectionString()

	courierStore, err := models.NewPgCourierStore(context.Background(), connStr)
	testutil.AssertNoErr(t, err)

	locationStore, err := models.NewPgLocationStore(context.Background(), connStr)
	testutil.AssertNoErr(t, err)

	courierEventHandler := handlers.NewCourierEventHandler(courierStore, locationStore)

	t.Run("creates new courier and initial location", func(t *testing.T) {
		wantCourier := testdata.VolenCourier
		wantLocation := models.Location{CourierID: wantCourier.ID}

		payload := svcevents.CourierCreatedEvent{
			ID:   wantCourier.ID,
			Name: wantCourier.Name,
		}
		event := events.NewTypedEvent(svcevents.COURIER_CREATED_EVENT_ID, wantCourier.ID, payload)

		err := courierEventHandler.HandleCourierCreatedEvent(event)
		testutil.AssertNoErr(t, err)

		gotCourier, err := courierStore.GetCourierByID(wantCourier.ID)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, gotCourier, wantCourier)

		gotLocation, err := locationStore.GetLocationByCourierID(wantCourier.ID)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, gotLocation, wantLocation)

	})

	t.Run("deletes courier and associated location", func(t *testing.T) {
		want := testdata.VolenCourier

		payload := svcevents.CourierDeletedEvent{
			ID: want.ID,
		}
		event := events.NewTypedEvent(svcevents.COURIER_DELETED_EVENT_ID, want.ID, payload)

		err := courierEventHandler.HandleCourierDeletedEvent(event)
		testutil.AssertNoErr(t, err)

		_, err = courierStore.GetCourierByID(want.ID)

		testutil.AssertError(t, err, storeerrors.ErrNotFound)

		_, err = locationStore.GetLocationByCourierID(want.ID)

		testutil.AssertError(t, err, storeerrors.ErrNotFound)
	})
}
