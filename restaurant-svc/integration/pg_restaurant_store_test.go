package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/restaurant-svc/handlers"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
	td "github.com/VitoNaychev/food-app/restaurant-svc/testdata"
	"github.com/VitoNaychev/food-app/restaurant-svc/testutil"
)

func DummyHander(w http.ResponseWriter, r *http.Request) {}

func TestCustomerServerOperations(t *testing.T) {
	connStr := SetupDatabaseContainer(t)

	restaurantStore, err := models.NewPgRestaurantStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	restaurantServer := handlers.NewRestaurantServer(env.SecretKey, env.ExpiresAt, &restaurantStore)

	server := handlers.NewRouterServer(restaurantServer,
		http.HandlerFunc(DummyHander),
		http.HandlerFunc(DummyHander),
		http.HandlerFunc(DummyHander))

	var shackJWT string

	t.Run("create new restaurant", func(t *testing.T) {
		request := handlers.NewCreateRestaurantRequest(td.ShackRestaurant)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		wantRestaurant := handlers.RestaurantToRestaurantResponse(td.ShackRestaurant)

		var got handlers.CreateRestaurantResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got.Restaurant, wantRestaurant)

		shackJWT = got.JWT.Token
	})

	t.Run("get restaurant", func(t *testing.T) {
		request := handlers.NewGetRestaurantRequest(shackJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		wantRestaurant := handlers.RestaurantToRestaurantResponse(td.ShackRestaurant)

		var got handlers.RestaurantResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, wantRestaurant)
	})

	t.Run("udpate restaurant", func(t *testing.T) {
		updateRestaurant := td.ShackRestaurant
		updateRestaurant.Name = "Chicken Snack"

		request := handlers.NewUpdateRestaurantRequest(shackJWT, updateRestaurant)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		wantRestaurant := handlers.RestaurantToRestaurantResponse(updateRestaurant)

		var got handlers.RestaurantResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, wantRestaurant)
	})

	t.Run("delete restaurant", func(t *testing.T) {
		request := handlers.NewDeleteRestaruantRequest(shackJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("get deleted restaurant", func(t *testing.T) {
		request := handlers.NewGetRestaurantRequest(shackJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
	})
}
