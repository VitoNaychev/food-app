package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/courier-svc/handlers"
	"github.com/VitoNaychev/food-app/courier-svc/models"
	td "github.com/VitoNaychev/food-app/courier-svc/testdata"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/integrationutil"
	"github.com/VitoNaychev/food-app/parser"
	"github.com/VitoNaychev/food-app/pgconfig"
	"github.com/VitoNaychev/food-app/testutil"
)

type DummyEventPublisher struct{}

func (d *DummyEventPublisher) Publish(string, events.InterfaceEvent) error {
	return nil
}

func TestCustomerServerOperations(t *testing.T) {
	config := pgconfig.GetConfigFromEnv(env)
	integrationutil.SetupDatabaseContainer(t, &config, "../sql-scripts/init.sql")

	courierStore, err := models.NewPgCourierStore(context.Background(), config.GetConnectionString())
	if err != nil {
		t.Fatal(err)
	}

	server := handlers.NewCourierServer(env.SecretKey, env.ExpiresAt, &courierStore, &DummyEventPublisher{})

	var michaelJWT string

	t.Run("create new courier", func(t *testing.T) {
		request := handlers.NewCreateCourierRequest(td.MichaelCourier)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		wantCourier := handlers.CourierToCourierResponse(td.MichaelCourier)

		got := parser.FromJSON[handlers.CreateCourierResponse](response.Body)

		testutil.AssertEqual(t, got.Courier, wantCourier)

		michaelJWT = got.JWT.Token
	})

	t.Run("get courier", func(t *testing.T) {
		request := handlers.NewGetCourierRequest(michaelJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		wantCourier := handlers.CourierToCourierResponse(td.MichaelCourier)

		got := parser.FromJSON[handlers.CourierResponse](response.Body)

		testutil.AssertEqual(t, got, wantCourier)
	})

	t.Run("udpate courier", func(t *testing.T) {
		updateCourier := td.MichaelCourier
		updateCourier.LastName = "Scarn"

		request := handlers.NewUpdateCourierRequest(michaelJWT, updateCourier)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		wantCourier := handlers.CourierToCourierResponse(updateCourier)

		got := parser.FromJSON[handlers.CourierResponse](response.Body)

		testutil.AssertEqual(t, got, wantCourier)
	})

	t.Run("delete courier", func(t *testing.T) {
		request := handlers.NewDeleteCourierRequest(michaelJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("get deleted courier", func(t *testing.T) {
		request := handlers.NewGetCourierRequest(michaelJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
	})
}
