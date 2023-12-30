package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/VitoNaychev/food-app/integrationutil"
	"github.com/VitoNaychev/food-app/parser"
	"github.com/VitoNaychev/food-app/pgconfig"
	"github.com/VitoNaychev/food-app/restaurant-svc/handlers"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
	td "github.com/VitoNaychev/food-app/restaurant-svc/testdata"
	"github.com/VitoNaychev/food-app/testutil"
)

func TestHoursServerOperations(t *testing.T) {
	config := pgconfig.GetConfigFromEnv(env)
	integrationutil.SetupDatabaseContainer(t, &config)

	connStr := config.GetConnectionString()

	hoursStore, err := models.NewPgHoursStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	restaurantStore, err := models.NewPgRestaurantStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	restaurantServer := handlers.NewRestaurantServer(env.SecretKey, env.ExpiresAt, &restaurantStore)
	hoursServer := handlers.NewHoursServer(env.SecretKey, &hoursStore, &restaurantStore)

	server := handlers.NewRouterServer(restaurantServer, DummyHandler, hoursServer, DummyHandler)

	shackJWT, err := createRestaurant(server, td.ShackRestaurant)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("creates hours", func(t *testing.T) {
		request := handlers.NewCreateHoursRequest(shackJWT, td.ShackHours)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := handlers.HoursArrToHoursResponseArr(td.ShackHours)
		got := parser.FromJSON[[]handlers.HoursResponse](response.Body)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("get hours", func(t *testing.T) {
		request := handlers.NewGetHoursRequest(shackJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := handlers.HoursArrToHoursResponseArr(td.ShackHours)
		got := parser.FromJSON[[]handlers.HoursResponse](response.Body)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("update hours", func(t *testing.T) {
		updateHours := make([]models.Hours, 2)
		copy(updateHours, td.ShackHours[4:6])
		updateHours[0].Opening, _ = time.Parse("15:04", "13:00")
		updateHours[1].Opening, _ = time.Parse("15:04", "13:00")
		request := handlers.NewUpdateHoursRequest(shackJWT, updateHours)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := handlers.HoursArrToHoursResponseArr(updateHours)
		got := parser.FromJSON[[]handlers.HoursResponse](response.Body)

		testutil.AssertEqual(t, got, want)
	})
}
