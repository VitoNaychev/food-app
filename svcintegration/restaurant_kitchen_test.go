package svcintegration

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/integrationutil"
	"github.com/VitoNaychev/food-app/restaurant-svc/handlers"
	"github.com/VitoNaychev/food-app/restaurant-svc/testdata"
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

	t.Run("kitchen-svc repicates new restaurant from restaurant-svc", func(t *testing.T) {
		request := handlers.NewCreateRestaurantRequest(testdata.ShackRestaurant)
		response := httptest.NewRecorder()

		restaurantService.router.ServeHTTP(response, request)
		time.Sleep(time.Second)

		got, err := kitchenService.restaurantStore.GetRestaurantByID(testdata.ShackRestaurant.ID)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, testdata.ShackRestaurant.ID, got.ID)
	})

}
