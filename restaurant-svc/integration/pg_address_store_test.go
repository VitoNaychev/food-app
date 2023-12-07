package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/restaurant-svc/handlers"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
	"github.com/VitoNaychev/food-app/restaurant-svc/testdata"
	td "github.com/VitoNaychev/food-app/restaurant-svc/testdata"
	"github.com/VitoNaychev/food-app/restaurant-svc/testutil"
)

func TestAddressServerOperations(t *testing.T) {
	connStr := SetupDatabaseContainer(t)

	addressStore, err := models.NewPgAddressStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	restaurantStore, err := models.NewPgRestaurantStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	restaurantServer := handlers.NewRestaurantServer(env.SecretKey, env.ExpiresAt, &restaurantStore)
	addressServer := handlers.NewAddressServer(env.SecretKey, &addressStore, &restaurantStore)

	server := handlers.NewRouterServer(restaurantServer, addressServer, DummyHandler, DummyHandler)

	shackJWT, err := createRestaurant(server, td.ShackRestaurant)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("creates address", func(t *testing.T) {
		request := handlers.NewCreateAddressRequest(shackJWT, td.ShackAddress)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		var got models.Address
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, td.ShackAddress)
	})

	t.Run("gets address", func(t *testing.T) {
		request := handlers.NewGetAddressRequest(shackJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		var got models.Address
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, td.ShackAddress)
	})

	t.Run("updates address", func(t *testing.T) {
		updateRestaurant := testdata.ShackAddress
		updateRestaurant.City = "Varna"

		request := handlers.NewUpdateAddressRequest(shackJWT, updateRestaurant)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		var got models.Address
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, updateRestaurant)
	})

}
