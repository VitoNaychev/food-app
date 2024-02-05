package handlers_test

import (
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

func TestLocationRequestValidation(t *testing.T) {
	courierStore := &stubs.StubCourierStore{
		Couriers: []models.Courier{testdata.VolenCourier},
	}
	locationStore := &stubs.StubLocationStore{
		Locations: []models.Location{testdata.VolenLocation},
	}

	server := handlers.NewLocationServer(env.SecretKey, locationStore, courierStore)

	volenJWT, _ := auth.GenerateJWT(env.SecretKey, env.ExpiresAt, testdata.VolenCourier.ID)

	cases := map[string]*http.Request{
		"update location": handlers.NewUpdateLocationRequest(volenJWT, 361, 181),
	}

	tabletests.RunRequestValidationTests(t, server, cases)
}

func TestUpdateCourierLocation(t *testing.T) {
	courierStore := &stubs.StubCourierStore{
		Couriers: []models.Courier{testdata.VolenCourier},
	}

	locationStore := &stubs.StubLocationStore{
		Locations: []models.Location{testdata.VolenLocation},
	}

	server := handlers.NewLocationServer(env.SecretKey, locationStore, courierStore)

	volenJWT, _ := auth.GenerateJWT(env.SecretKey, env.ExpiresAt, testdata.VolenCourier.ID)

	t.Run("updates courier location", func(t *testing.T) {
		want := models.Location{
			CourierID: testdata.VolenCourier.ID,
			Lat:       42.6492518454,
			Lon:       23.3450012782,
		}

		request := handlers.NewUpdateLocationRequest(volenJWT, want.Lat, want.Lon)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertEqual(t, locationStore.UpdatedLocation, want)

		got, err := validation.ValidateBody[models.Location](response.Body)
		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, got, want)
	})
}

func TestGetLocation(t *testing.T) {
	courierStore := &stubs.StubCourierStore{
		Couriers: []models.Courier{testdata.VolenCourier},
	}

	locationStore := &stubs.StubLocationStore{
		Locations: []models.Location{testdata.VolenLocation},
	}

	server := handlers.NewLocationServer(env.SecretKey, locationStore, courierStore)

	volenJWT, _ := auth.GenerateJWT(env.SecretKey, env.ExpiresAt, testdata.VolenCourier.ID)

	t.Run("gets courier location", func(t *testing.T) {
		want := handlers.LocationToGetLocationResponse(testdata.VolenLocation)

		request, _ := http.NewRequest(http.MethodGet, "/delivery/location/", nil)
		response := httptest.NewRecorder()

		request.Header.Add("Token", volenJWT)

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[handlers.GetLocationResponse](response.Body)
		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, got, want)
	})
}
