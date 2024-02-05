package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/delivery-svc/handlers"
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/delivery-svc/stubs"
	"github.com/VitoNaychev/food-app/delivery-svc/testdata"
	"github.com/VitoNaychev/food-app/testutil"
	"github.com/VitoNaychev/food-app/testutil/tabletests"
	"github.com/VitoNaychev/food-app/validation"
)

func TestDeliveryEndpointAuthentication(t *testing.T) {
	courierStore := &stubs.StubCourierStore{}

	server := handlers.NewDeliveryServer(env.SecretKey, nil, nil, courierStore)

	invalidJWT := "invalidJWT"
	cases := map[string]*http.Request{
		"change ticket state": handlers.NewChangeDeliveryStateRequest(invalidJWT, models.PICKUP_DELIVERY),
	}

	tabletests.RunAuthenticationTests(t, server, cases)
}

func TestDeliveryRequestValidation(t *testing.T) {
	courierStore := &stubs.StubCourierStore{
		Couriers: []models.Courier{testdata.VolenCourier},
	}

	deliveryStore := &stubs.StubDeliveryStore{
		Deliveries: []models.Delivery{testdata.VolenDelivery},
	}

	addressStore := &stubs.StubAddressStore{
		Addresses: []models.Address{testdata.VolenPickupAddress, testdata.VolenDeliveryAddress},
	}

	server := handlers.NewDeliveryServer(env.SecretKey, deliveryStore, addressStore, courierStore)

	volenJWT, _ := auth.GenerateJWT(env.SecretKey, env.ExpiresAt, testdata.VolenCourier.ID)

	cases := map[string]*http.Request{
		"change ticket state": handlers.NewChangeDeliveryStateRequest(volenJWT, -1),
	}

	tabletests.RunRequestValidationTests(t, server, cases)
}

func TestDeliveryStateTransitions(t *testing.T) {
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

	server := handlers.NewDeliveryServer(env.SecretKey, deliveryStore, nil, courierStore)

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

	t.Run("returns delivery status on state transition request", func(t *testing.T) {
		aliceJWT, _ := auth.GenerateJWT(env.SecretKey, env.ExpiresAt, testdata.AliceCourier.ID)

		state, _ := models.StateValueToStateName(models.ON_ROUTE)
		want := handlers.DeliveryStateTransitionResponse{
			ID:    testdata.AliceDelivery.ID,
			State: state,
		}

		request := handlers.NewChangeDeliveryStateRequest(aliceJWT, models.PICKUP_DELIVERY)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[handlers.DeliveryStateTransitionResponse](response.Body)
		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, got, want)
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
	})
}

func TestGetDelivery(t *testing.T) {
	courierStore := &stubs.StubCourierStore{
		Couriers: []models.Courier{testdata.VolenCourier, testdata.PeterCourier},
	}

	deliveryStore := &stubs.StubDeliveryStore{
		Deliveries: []models.Delivery{testdata.VolenDelivery},
	}

	addressStore := &stubs.StubAddressStore{
		Addresses: []models.Address{testdata.VolenPickupAddress, testdata.VolenDeliveryAddress},
	}

	server := handlers.NewDeliveryServer(env.SecretKey, deliveryStore, addressStore, courierStore)

	t.Run("returns current delivery info on GET", func(t *testing.T) {
		want := handlers.NewGetDeliveryResponse(testdata.VolenDelivery, testdata.VolenPickupAddress, testdata.VolenDeliveryAddress)

		volenJWT, _ := auth.GenerateJWT(env.SecretKey, env.ExpiresAt, testdata.VolenCourier.ID)

		request, _ := http.NewRequest(http.MethodGet, "/delivery/", nil)
		response := httptest.NewRecorder()

		request.Header.Add("Token", volenJWT)

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[handlers.GetDeliveryResponse](response.Body)
		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns empty body on no active deliveries", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(env.SecretKey, env.ExpiresAt, testdata.PeterCourier.ID)

		request, _ := http.NewRequest(http.MethodGet, "/delivery/", nil)
		response := httptest.NewRecorder()

		request.Header.Add("Token", peterJWT)

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		AssertEmptyBody(t, response.Body)
	})
}

func AssertEmptyBody(t testing.TB, body io.Reader) {
	t.Helper()
	var maxRequestSize int64 = 10000
	content, err := io.ReadAll(io.LimitReader(body, maxRequestSize))

	if err != nil {
		t.Errorf("error reading response body: %v", err)
	}

	if len(content) != 0 {
		t.Errorf("response body is not empty")
	}
}
