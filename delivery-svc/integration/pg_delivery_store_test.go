package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/delivery-svc/handlers"
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/delivery-svc/testdata"
	"github.com/VitoNaychev/food-app/integrationutil"
	"github.com/VitoNaychev/food-app/pgconfig"
	"github.com/VitoNaychev/food-app/testutil"
	"github.com/VitoNaychev/food-app/validation"
)

func TestDeliveryHandlerOperations(t *testing.T) {
	config := pgconfig.GetConfigFromEnv(env)
	integrationutil.SetupDatabaseContainer(t, &config, "../sql-scripts/init.sql")

	connStr := config.GetConnectionString()

	courierStore, err := models.NewPgCourierStore(context.Background(), connStr)
	testutil.AssertNoErr(t, err)
	initCouriersTable(t, courierStore)

	addressStore, err := models.NewPgAddressStore(context.Background(), connStr)
	testutil.AssertNoErr(t, err)
	initAddressesTable(t, addressStore)

	deliveryStore, err := models.NewPgDeliveryStore(context.Background(), connStr)
	testutil.AssertNoErr(t, err)
	initDeliveriesTable(t, deliveryStore)

	server := handlers.NewDeliveryServer(env.SecretKey, deliveryStore, addressStore, courierStore)

	volenJWT, _ := auth.GenerateJWT(env.SecretKey, env.ExpiresAt, testdata.VolenCourier.ID)

	t.Run("gets courier's active delivery", func(t *testing.T) {
		want := handlers.NewGetDeliveryResponse(testdata.VolenActiveDelivery, testdata.VolenPickupAddress, testdata.VolenDeliveryAddress)

		request := handlers.NewGetActiveDeliveryRequest(volenJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[handlers.GetDeliveryResponse](response.Body)
		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, got, want)
	})

	t.Run("updates delivery state", func(t *testing.T) {
		request := handlers.NewChangeDeliveryStateRequest(volenJWT, models.PICKUP_DELIVERY)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		stateName, _ := models.StateValueToStateName(models.ON_ROUTE)
		want := handlers.DeliveryStateTransitionResponse{
			ID:    testdata.VolenActiveDelivery.ID,
			State: stateName,
		}
		got, err := validation.ValidateBody[handlers.DeliveryStateTransitionResponse](response.Body)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, got, want)
	})
}
