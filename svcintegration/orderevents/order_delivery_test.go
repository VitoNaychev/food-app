package orderevents

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/integrationutil"
	"github.com/VitoNaychev/food-app/order-svc/handlers"
	"github.com/VitoNaychev/food-app/svcintegration/services"
	"github.com/VitoNaychev/food-app/testutil"
)

func TestDeliveryServiceOrderEventHandlers(t *testing.T) {
	_, brokersAddrs := integrationutil.SetupKafkaContainer(t)

	orderKeys := []string{}
	orderEnv, err := appenv.LoadEnviornment("../../order-svc/test.env", orderKeys)
	if err != nil {
		t.Fatalf("Failed to load order env: %v", err)
	}
	orderEnv.KafkaBrokers = brokersAddrs

	deliveryKeys := []string{}
	deliveryEnv, err := appenv.LoadEnviornment("../../delivery-svc/test.env", deliveryKeys)
	if err != nil {
		t.Fatalf("Failed to load delivery env: %v", err)
	}
	deliveryEnv.KafkaBrokers = brokersAddrs

	orderService := services.SetupOrderService(t, orderEnv, ":4040")
	orderService.Run()
	defer orderService.Stop()

	deliveryService := services.SetupDeliveryService(t, deliveryEnv, ":8080")
	deliveryService.Run()
	defer deliveryService.Stop()

	initDeliveryServiceTables(t, deliveryService)

	t.Run("delivery-svc creates coresponding delivery on ORDER_CREATED_EVENT", func(t *testing.T) {
		createOrderRequestBody := handlers.NewCeateOrderRequestBody(peterCreatedOrder, peterCreatedOrderItems, chickenShackAddress, peterAddress1)
		request := handlers.NewCreateOrderRequest(strconv.Itoa(peterCustomerID), createOrderRequestBody)
		response := httptest.NewRecorder()

		orderService.Server.Handler.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		time.Sleep(time.Second)

		gotDelivery, err := deliveryService.DeliveryStore.GetDeliveryByID(peterCreatedOrder.ID)
		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, gotDelivery, volenDelivery)

		gotPickupAddress, err := deliveryService.AddressStore.GetAddressByID(gotDelivery.PickupAddressID)
		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, gotPickupAddress, volenPickupAddress)

		gotDeliveryAddress, err := deliveryService.AddressStore.GetAddressByID(gotDelivery.DeliveryAddressID)
		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, gotDeliveryAddress, volenDeliveryAddress)
	})
}

func initDeliveryServiceTables(t testing.TB, svc services.DeliveryService) {
	testutil.AssertNoErr(t, svc.CourierStore.CreateCourier(&volenCourier))
}
