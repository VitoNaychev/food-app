package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/integrationutil"
	"github.com/VitoNaychev/food-app/parser"
	"github.com/VitoNaychev/food-app/restaurant-svc/handlers"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
	"github.com/VitoNaychev/food-app/restaurant-svc/testdata"
	td "github.com/VitoNaychev/food-app/restaurant-svc/testdata"
	"github.com/VitoNaychev/food-app/testutil"
)

type FakePublisher struct {
}

func (s *FakePublisher) Publish(topic string, event events.Event) error {
	return nil
}

func TestMenuServerOperations(t *testing.T) {
	connStr := integrationutil.SetupDatabaseContainer(t, env)

	menuStore, err := models.NewPgMenuStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	hoursStore, err := models.NewPgHoursStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

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
	hoursServer := handlers.NewHoursServer(env.SecretKey, &hoursStore, &restaurantStore)
	menuServer := handlers.NewMenuServer(env.SecretKey, &menuStore, &restaurantStore, &FakePublisher{})

	server := handlers.NewRouterServer(restaurantServer, addressServer, hoursServer, menuServer)

	dominosJWT, err := createRestaurant(server, td.DominosRestaurant)
	if err != nil {
		t.Fatal(err)
	}

	err = createAddress(server, dominosJWT, testdata.DominosAddress)
	if err != nil {
		t.Fatal(err)
	}

	err = createHours(server, dominosJWT, testdata.DominosHours)
	if err != nil {
		t.Fatal(err)
	}

	// In the testdata package Dominos is the second restaurant, but
	// in this case it is the first and only restaurant created, so
	// change the id of the menu item that we test to reflect that.
	testItem := td.DominosMenu[0]
	testItem.RestaurantID = 1

	t.Run("creates menu item", func(t *testing.T) {
		request := handlers.NewCreateMenuItemRequest(dominosJWT, testItem)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got := parser.FromJSON[models.MenuItem](response.Body)

		testutil.AssertEqual(t, got, testItem)
	})

	t.Run("gets restaurant menu", func(t *testing.T) {
		request := handlers.NewGetMenuRequest(dominosJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := []models.MenuItem{testItem}
		got := parser.FromJSON[[]models.MenuItem](response.Body)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("updates menu item", func(t *testing.T) {
		updateItem := testItem
		updateItem.Name = "Master Burger Pizza"
		updateItem.Details = "Bestest pizza in the world, bruh"

		request := handlers.NewUpdateMenuItemRequest(dominosJWT, updateItem)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got := parser.FromJSON[models.MenuItem](response.Body)

		testutil.AssertEqual(t, got, updateItem)
	})

	t.Run("deletes menu item", func(t *testing.T) {
		request := handlers.NewDeleteMenuItemRequest(dominosJWT, handlers.DeleteMenuItemRequest{ID: 1})
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("get menu after delete", func(t *testing.T) {
		request := handlers.NewGetMenuRequest(dominosJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := []models.MenuItem{}
		got := parser.FromJSON[[]models.MenuItem](response.Body)

		testutil.AssertEqual(t, got, want)
	})
}
