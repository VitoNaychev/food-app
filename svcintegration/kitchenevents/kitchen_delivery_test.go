package kitchenevents

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/integrationutil"
	"github.com/VitoNaychev/food-app/kitchen-svc/handlers"
	"github.com/VitoNaychev/food-app/svcintegration/services"
	"github.com/VitoNaychev/food-app/testutil"
)

func TestKitchenDomainEvents(t *testing.T) {
	_, brokersAddrs := integrationutil.SetupKafkaContainer(t)

	kitchenKeys := []string{"SECRET", "EXPIRES_AT"}
	kitchenEnv, err := appenv.LoadEnviornment("../../kitchen-svc/test.env", kitchenKeys)
	if err != nil {
		t.Fatalf("Failed to load kitchen env: %v", err)
	}
	kitchenEnv.KafkaBrokers = brokersAddrs

	deliveryKeys := []string{}
	deliveryEnv, err := appenv.LoadEnviornment("../../delivery-svc/test.env", deliveryKeys)
	if err != nil {
		t.Fatalf("Failed to load delivery env: %v", err)
	}
	deliveryEnv.KafkaBrokers = brokersAddrs

	kitchenService := services.SetupKitchenService(t, kitchenEnv, ":9090")
	kitchenService.Run()
	defer kitchenService.Stop()

	deliveryService := services.SetupDeliveryService(t, deliveryEnv, ":8080")
	deliveryService.Run()
	defer deliveryService.Stop()

	initKitchenServiceTables(t, kitchenService)
	initDeliveryServiceTables(t, deliveryService)

	shackJWT, _ := auth.GenerateJWT(kitchenEnv.SecretKey, kitchenEnv.ExpiresAt, shackRestaurant.ID)

	t.Run("delivery-svc updates delivery state to IN_PROGRESS when kitchen begins preparing ticket", func(t *testing.T) {
		readyByStr := "23:59"
		readyByTime, _ := handlers.ParseTimeAndSetDate(readyByStr)

		request := handlers.NewBeginPreparingTicketRequest(shackJWT, shackTicket.ID, readyByStr)
		response := httptest.NewRecorder()

		kitchenService.Server.Handler.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		time.Sleep(time.Second)

		wantState := models.IN_PROGRESS
		wantReadyBy := readyByTime

		got, err := deliveryService.DeliveryStore.GetDeliveryByID(volenDelivery.ID)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, got.State, wantState)
		testutil.AssertEqual(t, got.ReadyBy, wantReadyBy)
	})

	t.Run("delivery-svc updates delivery state to READY_FOR_PICKUP when kitchen finishes preparing ticket", func(t *testing.T) {
		request := handlers.NewFinishPreparingTicketRequest(shackJWT, shackTicket.ID)
		response := httptest.NewRecorder()

		kitchenService.Server.Handler.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		time.Sleep(time.Second)

		wantState := models.READY_FOR_PICKUP

		got, err := deliveryService.DeliveryStore.GetDeliveryByID(volenDelivery.ID)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, got.State, wantState)
	})
}

func initKitchenServiceTables(t testing.TB, svc services.KitchenService) {
	testutil.AssertNoErr(t, svc.RestaurantStore.CreateRestaurant(&shackRestaurant))
	testutil.AssertNoErr(t, svc.MenuItemStore.CreateMenuItem(&shackMenuItem))
	testutil.AssertNoErr(t, svc.TicketStore.CreateTicket(&shackTicket))
	testutil.AssertNoErr(t, svc.TicketItemStore.CreateTicketItem(&shackTicketItem))
}

func initDeliveryServiceTables(t testing.TB, svc services.DeliveryService) {
	testutil.AssertNoErr(t, svc.AddressStore.CreateAddress(&volenPickupAddress))
	testutil.AssertNoErr(t, svc.AddressStore.CreateAddress(&volenDeliveryAddress))
	testutil.AssertNoErr(t, svc.CourierStore.CreateCourier(&volenCourier))
	testutil.AssertNoErr(t, svc.DeliveryStore.CreateDelivery(&volenDelivery))
}
