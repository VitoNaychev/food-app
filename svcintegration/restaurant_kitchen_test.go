package svcintegration

import (
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
	"github.com/VitoNaychev/food-app/testutil"
)

var restaurantKeys = []string{"SECRET", "EXPIRES_AT"}
var kitchenKeys = []string{}

func TestRestaurantDomainEvents(t *testing.T) {
	_, brokersAddrs := integrationutil.SetupKafkaContainer(t)

	restaurantEnv, err := appenv.LoadEnviornment("../restaurant-svc/test.env", restaurantKeys)
	if err != nil {
		t.Fatalf("Failed to load restaurant env: %v", err)
	}
	restaurantEnv.KafkaBrokers = brokersAddrs

	kitchenEnv, err := appenv.LoadEnviornment("../kitchen-svc/test.env", kitchenKeys)
	if err != nil {
		t.Fatalf("Failed to load restaurant env: %v", err)
	}
	kitchenEnv.KafkaBrokers = brokersAddrs

	restaurantService := SetupRestaurantService(t, restaurantEnv, ":8080")
	restaurantService.Run()
	defer restaurantService.Stop()

	kitchenService := SetupKitchenService(t, kitchenEnv, ":9090")
	kitchenService.Run()
	defer kitchenService.Stop()

	shackJWT := ""

	t.Run("kitchen-svc repicates new restaurant", func(t *testing.T) {
		request := handlers.NewCreateRestaurantRequest(testdata.ShackRestaurant)
		response := httptest.NewRecorder()

		restaurantService.router.ServeHTTP(response, request)
		testutil.AssertStatus(t, response.Code, http.StatusOK)

		shackJWT = getJWTFromResponseBody(response.Body)

		time.Sleep(time.Second)

		got, err := kitchenService.restaurantStore.GetRestaurantByID(testdata.ShackRestaurant.ID)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, got.ID, testdata.ShackRestaurant.ID)
	})

	{
		request := handlers.NewCreateAddressRequest(shackJWT, testdata.ShackAddress)
		response := httptest.NewRecorder()

		restaurantService.router.ServeHTTP(response, request)
		testutil.AssertStatus(t, response.Code, http.StatusOK)
	}

	{
		request := handlers.NewCreateHoursRequest(shackJWT, testdata.ShackHours)
		response := httptest.NewRecorder()

		restaurantService.router.ServeHTTP(response, request)
		testutil.AssertStatus(t, response.Code, http.StatusOK)
	}

	t.Run("kitchen-svc replicates new menu item", func(t *testing.T) {
		request := handlers.NewCreateMenuItemRequest(shackJWT, testdata.ShackMenu[0])
		response := httptest.NewRecorder()

		restaurantService.router.ServeHTTP(response, request)
		testutil.AssertStatus(t, response.Code, http.StatusOK)

		time.Sleep(time.Second)

		got, err := kitchenService.menuItemStore.GetMenuItemByID(testdata.ShackMenu[0].ID)

		testutil.AssertNoErr(t, err)
		AssertMenuItem(t, got, testdata.ShackMenu[0])
	})

	t.Run("kitchen-svc updates menu item", func(t *testing.T) {
		menuItem := testdata.ShackMenu[0]
		menuItem.Name = "XXXXL Duner"

		request := handlers.NewUpdateMenuItemRequest(shackJWT, menuItem)
		response := httptest.NewRecorder()

		restaurantService.router.ServeHTTP(response, request)
		testutil.AssertStatus(t, response.Code, http.StatusOK)

		time.Sleep(time.Second)

		got, err := kitchenService.menuItemStore.GetMenuItemByID(menuItem.ID)

		testutil.AssertNoErr(t, err)
		AssertMenuItem(t, got, menuItem)
	})

	t.Run("kitchen-svc deletes menu item", func(t *testing.T) {
		wantID := testdata.ShackMenu[0].ID

		request := handlers.NewDeleteMenuItemRequest(shackJWT, handlers.DeleteMenuItemRequest{ID: wantID})
		response := httptest.NewRecorder()

		restaurantService.router.ServeHTTP(response, request)
		testutil.AssertStatus(t, response.Code, http.StatusOK)

		time.Sleep(time.Second)

		_, err := kitchenService.menuItemStore.GetMenuItemByID(wantID)
		AssertError(t, err, storeerrors.ErrNotFound)
	})

	{
		request := handlers.NewCreateMenuItemRequest(shackJWT, testdata.ShackMenu[0])
		response := httptest.NewRecorder()

		restaurantService.router.ServeHTTP(response, request)
		testutil.AssertStatus(t, response.Code, http.StatusOK)
	}

	t.Run("kitchen-svc deletes restaurant and associated menu items", func(t *testing.T) {
		request := handlers.NewDeleteRestaurantRequest(shackJWT)
		response := httptest.NewRecorder()

		restaurantService.router.ServeHTTP(response, request)
		testutil.AssertStatus(t, response.Code, http.StatusOK)

		time.Sleep(time.Second)

		_, err := kitchenService.restaurantStore.GetRestaurantByID(testdata.ShackRestaurant.ID)
		AssertError(t, err, storeerrors.ErrNotFound)

		_, err = kitchenService.menuItemStore.GetMenuItemByID(testdata.ShackMenu[0].ID)
		AssertError(t, err, storeerrors.ErrNotFound)
	})
}

func AssertError(t testing.TB, got error, want error) {
	t.Helper()

	if got == nil {
		t.Fatalf("expected error but didn't get one")
	}

	if got != want {
		t.Errorf("got error %v, want %v", got, want)
	}
}

func AssertMenuItem(t testing.TB, got kitchenmodels.MenuItem, want restaurantmodels.MenuItem) {
	t.Helper()

	testutil.AssertEqual(t, got.ID, want.ID)
	testutil.AssertEqual(t, got.Name, want.Name)
	testutil.AssertEqual(t, got.Price, want.Price)
	testutil.AssertEqual(t, got.RestaurantID, want.RestaurantID)
}
