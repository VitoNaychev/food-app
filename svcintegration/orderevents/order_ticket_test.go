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

func TestKitchenServiceOrderEventHandlers(t *testing.T) {
	_, brokersAddrs := integrationutil.SetupKafkaContainer(t)

	orderKeys := []string{}
	orderEnv, err := appenv.LoadEnviornment("../../order-svc/test.env", orderKeys)
	if err != nil {
		t.Fatalf("Failed to load order env: %v", err)
	}
	orderEnv.KafkaBrokers = brokersAddrs

	kitchenKeys := []string{}
	kitchenEnv, err := appenv.LoadEnviornment("../../kitchen-svc/test.env", kitchenKeys)
	if err != nil {
		t.Fatalf("Failed to load kitchen env: %v", err)
	}
	kitchenEnv.KafkaBrokers = brokersAddrs

	orderService := services.SetupOrderService(t, orderEnv, ":9090")
	orderService.Run()
	defer orderService.Stop()

	kitchenService := services.SetupKitchenService(t, kitchenEnv, ":8080")
	kitchenService.Run()
	defer kitchenService.Stop()

	t.Run("kitchen-svc creates coresponding ticket on ORDER_CREATED_EVENT", func(t *testing.T) {
		createOrderRequestBody := handlers.NewCeateOrderRequestBody(peterCreatedOrder, peterCreatedOrderItems, chickenShackAddress, peterAddress1)
		request := handlers.NewCreateOrderRequest(strconv.Itoa(peterCustomerID), createOrderRequestBody)
		response := httptest.NewRecorder()

		orderService.Server.Handler.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		time.Sleep(time.Second)

		gotTicket, err := kitchenService.TicketStore.GetTicketByID(peterCreatedOrder.ID)
		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, gotTicket, shackTicket)

		gotTicketItems, err := kitchenService.TicketItemStore.GetTicketItemsByTicketID(peterCreatedOrder.ID)
		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, gotTicketItems, shackTicketItems)
	})
}
