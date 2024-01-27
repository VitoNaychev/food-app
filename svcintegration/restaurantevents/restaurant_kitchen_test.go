package svcintegration

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/integrationutil"
	kitchenmodels "github.com/VitoNaychev/food-app/kitchen-svc/models"
	"github.com/VitoNaychev/food-app/restaurant-svc/handlers"
	restaurantmodels "github.com/VitoNaychev/food-app/restaurant-svc/models"
	"github.com/VitoNaychev/food-app/restaurant-svc/testdata"
	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/VitoNaychev/food-app/svcintegration/services"
	"github.com/VitoNaychev/food-app/testutil"
)

func TestRestaurantDomainEvents(t *testing.T) {
	_, brokersAddrs := integrationutil.SetupKafkaContainer(t)

	restaurantKeys := []string{"SECRET", "EXPIRES_AT"}
	restaurantEnv, err := appenv.LoadEnviornment("../../restaurant-svc/test.env", restaurantKeys)
	if err != nil {
		t.Fatalf("Failed to load restaurant env: %v", err)
	}
	restaurantEnv.KafkaBrokers = brokersAddrs

	kitchenKeys := []string{}
	kitchenEnv, err := appenv.LoadEnviornment("../../kitchen-svc/test.env", kitchenKeys)
	if err != nil {
		t.Fatalf("Failed to load kitchen env: %v", err)
	}
	kitchenEnv.KafkaBrokers = brokersAddrs

	restaurantService := services.SetupRestaurantService(t, restaurantEnv, ":8080")
	restaurantService.Run()
	defer restaurantService.Stop()

	kitchenService := services.SetupKitchenService(t, kitchenEnv, ":9090")
	kitchenService.Run()
	defer kitchenService.Stop()

	shackJWT := ""

	t.Run("kitchen-svc repicates new restaurant", func(t *testing.T) {
		request := handlers.NewCreateRestaurantRequest(testdata.ShackRestaurant)
		response := httptest.NewRecorder()

		restaurantService.Router.ServeHTTP(response, request)
		testutil.AssertStatus(t, response.Code, http.StatusOK)

		shackJWT = getJWTFromResponseBody(response.Body)

		time.Sleep(time.Second)

		got, err := kitchenService.RestaurantStore.GetRestaurantByID(testdata.ShackRestaurant.ID)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, got.ID, testdata.ShackRestaurant.ID)
	})

	{
		request := handlers.NewCreateAddressRequest(shackJWT, testdata.ShackAddress)
		response := httptest.NewRecorder()

		restaurantService.Router.ServeHTTP(response, request)
		testutil.AssertStatus(t, response.Code, http.StatusOK)
	}

	{
		request := handlers.NewCreateHoursRequest(shackJWT, testdata.ShackHours)
		response := httptest.NewRecorder()

		restaurantService.Router.ServeHTTP(response, request)
		testutil.AssertStatus(t, response.Code, http.StatusOK)
	}

	t.Run("kitchen-svc replicates new menu item", func(t *testing.T) {
		request := handlers.NewCreateMenuItemRequest(shackJWT, testdata.ShackMenu[0])
		response := httptest.NewRecorder()

		restaurantService.Router.ServeHTTP(response, request)
		testutil.AssertStatus(t, response.Code, http.StatusOK)

		time.Sleep(time.Second)

		got, err := kitchenService.MenuItemStore.GetMenuItemByID(testdata.ShackMenu[0].ID)

		testutil.AssertNoErr(t, err)
		AssertMenuItem(t, got, testdata.ShackMenu[0])
	})

	t.Run("kitchen-svc updates menu item", func(t *testing.T) {
		menuItem := testdata.ShackMenu[0]
		menuItem.Name = "XXXXL Duner"

		request := handlers.NewUpdateMenuItemRequest(shackJWT, menuItem)
		response := httptest.NewRecorder()

		restaurantService.Router.ServeHTTP(response, request)
		testutil.AssertStatus(t, response.Code, http.StatusOK)

		time.Sleep(time.Second)

		got, err := kitchenService.MenuItemStore.GetMenuItemByID(menuItem.ID)

		testutil.AssertNoErr(t, err)
		AssertMenuItem(t, got, menuItem)
	})

	t.Run("kitchen-svc deletes menu item", func(t *testing.T) {
		wantID := testdata.ShackMenu[0].ID

		request := handlers.NewDeleteMenuItemRequest(shackJWT, handlers.DeleteMenuItemRequest{ID: wantID})
		response := httptest.NewRecorder()

		restaurantService.Router.ServeHTTP(response, request)
		testutil.AssertStatus(t, response.Code, http.StatusOK)

		time.Sleep(time.Second)

		_, err := kitchenService.MenuItemStore.GetMenuItemByID(wantID)
		testutil.AssertError(t, err, storeerrors.ErrNotFound)
	})

	{
		request := handlers.NewCreateMenuItemRequest(shackJWT, testdata.ShackMenu[0])
		response := httptest.NewRecorder()

		restaurantService.Router.ServeHTTP(response, request)
		testutil.AssertStatus(t, response.Code, http.StatusOK)
	}

	t.Run("kitchen-svc deletes restaurant and associated menu items", func(t *testing.T) {
		request := handlers.NewDeleteRestaurantRequest(shackJWT)
		response := httptest.NewRecorder()

		restaurantService.Router.ServeHTTP(response, request)
		testutil.AssertStatus(t, response.Code, http.StatusOK)

		time.Sleep(time.Second)

		_, err := kitchenService.RestaurantStore.GetRestaurantByID(testdata.ShackRestaurant.ID)
		testutil.AssertError(t, err, storeerrors.ErrNotFound)

		_, err = kitchenService.MenuItemStore.GetMenuItemByID(testdata.ShackMenu[0].ID)
		testutil.AssertError(t, err, storeerrors.ErrNotFound)
	})
}

func AssertMenuItem(t testing.TB, got kitchenmodels.MenuItem, want restaurantmodels.MenuItem) {
	t.Helper()

	testutil.AssertEqual(t, got.ID, want.ID)
	testutil.AssertEqual(t, got.Name, want.Name)
	testutil.AssertEqual(t, got.Price, want.Price)
	testutil.AssertEqual(t, got.RestaurantID, want.RestaurantID)
}

func getJWTFromResponseBody(body io.Reader) string {
	var createRestaurantResponse handlers.CreateRestaurantResponse
	json.NewDecoder(body).Decode(&createRestaurantResponse)

	return createRestaurantResponse.JWT.Token
}
