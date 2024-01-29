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

	courierEventHandler := handlers.NewCourierEventHandler(courierStore)

	t.Run("creates new courier", func(t *testing.T) {
		want := testdata.VolenCourier

		payload := svcevents.CourierCreatedEvent{
			ID:   want.ID,
			Name: want.Name,
		}
		event := events.NewTypedEvent(svcevents.COURIER_CREATED_EVENT_ID, want.ID, payload)

		err := courierEventHandler.HandleCourierCreatedEvent(event)
		testutil.AssertNoErr(t, err)

		got, err := courierStore.GetCourierByID(want.ID)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, got, want)
	})

	t.Run("deletes courier", func(t *testing.T) {
		want := testdata.VolenCourier

		payload := svcevents.CourierDeletedEvent{
			ID: want.ID,
		}
		event := events.NewTypedEvent(svcevents.COURIER_DELETED_EVENT_ID, want.ID, payload)

		err := courierEventHandler.HandleCourierDeletedEvent(event)
		testutil.AssertNoErr(t, err)

		_, err = courierStore.GetCourierByID(want.ID)

		testutil.AssertError(t, err, storeerrors.ErrNotFound)
	})
}
