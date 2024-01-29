package integration

import (
	"context"
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

func TestLocationHandlerOperations(t *testing.T) {
	config := pgconfig.GetConfigFromEnv(env)
	integrationutil.SetupDatabaseContainer(t, &config, "../sql-scripts/init.sql")

	connStr := config.GetConnectionString()

	courierStore, err := models.NewPgCourierStore(context.Background(), connStr)
	testutil.AssertNoErr(t, err)
	initCouriersTable(t, courierStore)

	locationStore, err := models.NewPgLocationStore(context.Background(), connStr)
	testutil.AssertNoErr(t, err)
	initLocationsTable(t, locationStore)

	server := handlers.NewLocationServer(env.SecretKey, locationStore, courierStore)

	volenJWT, _ := auth.GenerateJWT(env.SecretKey, env.ExpiresAt, testdata.VolenCourier.ID)

	t.Run("gets courier location", func(t *testing.T) {
		want := testdata.VolenLocation

		request := handlers.NewGetLocationRequest(volenJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, response.Code)

		got, err := validation.ValidateBody[models.Location](response.Body)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, got, want)
	})

	t.Run("updates courier location", func(t *testing.T) {
		want := testdata.VolenLocation
		want.Lat = 42.69321155122
		want.Lon = 23.33248427579

		request := handlers.NewUpdateLocationRequest(volenJWT, want.Lat, want.Lon)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, response.Code)

		got, err := validation.ValidateBody[models.Location](response.Body)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, got, want)
	})
}
