package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/delivery-svc/handlers"
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/delivery-svc/stubs"
	"github.com/VitoNaychev/food-app/delivery-svc/testdata"
	"github.com/VitoNaychev/food-app/testutil"
	"github.com/VitoNaychev/food-app/testutil/tabletests"
)

var env appenv.Enviornment

func TestMain(m *testing.M) {
	keys := []string{"SECRET", "EXPIRES_AT"}

	var err error
	env, err = appenv.LoadEnviornment("../test.env", keys)
	if err != nil {
		testutil.HandleLoadEnviornmentError(err)
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}

func TestDeliveryEndpointAuthentication(t *testing.T) {
	courierStore := &stubs.StubCourierStore{}

	server := handlers.NewDeliveryServer(env.SecretKey, nil, courierStore)

	invalidJWT := "invalidJWT"
	cases := map[string]*http.Request{
		"change ticket state": handlers.NewChangeDeliveryStateRequest(invalidJWT, models.PICKUP_DELIVERY),
	}

	tabletests.RunAuthenticationTests(t, server, cases)
}

func TestDeliveryController(t *testing.T) {
	courierStore := &stubs.StubCourierStore{
		Couriers: []models.Courier{
			testdata.VolenCourier, testdata.PeterCourier, testdata.AliceCourier,
			testdata.JohnCourier, testdata.IvoCourier,
		},
	}

	deliveryStore := &stubs.StubDeliveryStore{
		Deliveries: []models.Delivery{
			testdata.VolenDelivery, testdata.PeterDelivery, testdata.AliceDelivery,
			testdata.JohnDelivery, testdata.IvoDelivery,
		},
	}

	server := handlers.NewDeliveryServer(env.SecretKey, deliveryStore, courierStore)

	t.Run("changes delivery state to ON_ROUTE on PICKUP_DELIVERY event", func(t *testing.T) {
		aliceJWT, _ := auth.GenerateJWT(env.SecretKey, env.ExpiresAt, testdata.AliceCourier.ID)
		want := testdata.AliceDelivery
		want.State = models.ON_ROUTE

		request := handlers.NewChangeDeliveryStateRequest(aliceJWT, models.PICKUP_DELIVERY)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertEqual(t, deliveryStore.UpdatedDelivery, want)
	})

	t.Run("changes delivery state to COMPLETED on COMPLETE_DELIVERY event", func(t *testing.T) {
		johnJWT, _ := auth.GenerateJWT(env.SecretKey, env.ExpiresAt, testdata.JohnCourier.ID)
		want := testdata.JohnDelivery
		want.State = models.COMPLETED

		request := handlers.NewChangeDeliveryStateRequest(johnJWT, models.COMPLETE_DELIVERY)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertEqual(t, deliveryStore.UpdatedDelivery, want)
	})

	t.Run("returns Bad Request if courier doesn't have an active delivery", func(t *testing.T) {
		ivoJWT, _ := auth.GenerateJWT(env.SecretKey, env.ExpiresAt, testdata.IvoCourier.ID)

		request := handlers.NewChangeDeliveryStateRequest(ivoJWT, models.COMPLETE_DELIVERY)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrNoActiveDeliveries)
		// testutil.AssertEqual(t, deliveryStore.UpdatedDelivery, want)
	})
}
