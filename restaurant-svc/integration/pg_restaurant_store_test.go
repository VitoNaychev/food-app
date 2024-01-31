package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/integrationutil"
	"github.com/VitoNaychev/food-app/parser"
	"github.com/VitoNaychev/food-app/pgconfig"
	"github.com/VitoNaychev/food-app/restaurant-svc/handlers"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
	td "github.com/VitoNaychev/food-app/restaurant-svc/testdata"
	"github.com/VitoNaychev/food-app/testutil"
	"github.com/VitoNaychev/food-app/testutil/dummies"
)

func TestCustomerServerOperations(t *testing.T) {
	config := pgconfig.GetConfigFromEnv(env)
	integrationutil.SetupDatabaseContainer(t, &config, "../sql-scripts/init.sql")

	connStr := config.GetConnectionString()

	restaurantStore, err := models.NewPgRestaurantStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	restaurantServer := handlers.NewRestaurantServer(env.SecretKey, env.ExpiresAt, &restaurantStore, &dummies.DummyPublisher{})

	server := handlers.NewRouterServer(restaurantServer,
		http.HandlerFunc(DummyHandler),
		http.HandlerFunc(DummyHandler),
		http.HandlerFunc(DummyHandler))

	var shackJWT string

	t.Run("create new restaurant", func(t *testing.T) {
		request := handlers.NewCreateRestaurantRequest(td.ShackRestaurant)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		wantRestaurant := handlers.RestaurantToRestaurantResponse(td.ShackRestaurant)

		got := parser.FromJSON[handlers.CreateRestaurantResponse](response.Body)

		testutil.AssertEqual(t, got.Restaurant, wantRestaurant)

		shackJWT = got.JWT.Token
	})

	t.Run("get restaurant", func(t *testing.T) {
		request := handlers.NewGetRestaurantRequest(shackJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		wantRestaurant := handlers.RestaurantToRestaurantResponse(td.ShackRestaurant)

		got := parser.FromJSON[handlers.RestaurantResponse](response.Body)

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

		got := parser.FromJSON[handlers.RestaurantResponse](response.Body)

		testutil.AssertEqual(t, got, wantRestaurant)
	})

	t.Run("delete restaurant", func(t *testing.T) {
		request := handlers.NewDeleteRestaurantRequest(shackJWT)
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
