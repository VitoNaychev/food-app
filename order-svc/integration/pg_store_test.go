package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/VitoNaychev/food-app/integrationutil"
	"github.com/VitoNaychev/food-app/order-svc/handlers"
	"github.com/VitoNaychev/food-app/order-svc/models"
	"github.com/VitoNaychev/food-app/order-svc/stubs"
	"github.com/VitoNaychev/food-app/order-svc/testdata"
	"github.com/VitoNaychev/food-app/pgconfig"
	"github.com/VitoNaychev/food-app/testutil"
	"github.com/VitoNaychev/food-app/testutil/dummies"
)

func TestOrderServerOperations(t *testing.T) {
	config := pgconfig.GetConfigFromEnv(env)
	integrationutil.SetupDatabaseContainer(t, &config, "../sql-scripts/init.sql")

	connStr := config.GetConnectionString()

	orderStore, err := models.NewPgOrderStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	orderItemStore, err := models.NewPgOrderItemStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	addressStore, err := models.NewPgAddressStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	server := handlers.NewOrderServer(orderStore, orderItemStore, addressStore, &dummies.DummyPublisher{}, stubs.StubVerifyJWT)

	peterJWT := strconv.Itoa(testdata.PeterCustomerID)
	createOrderRequestBody := handlers.NewCeateOrderRequestBody(testdata.PeterCreatedOrder, testdata.PeterCreatedOrderItems, testdata.ChickenShackAddress, testdata.PeterAddress1)

	request := handlers.NewCreateOrderRequest(peterJWT, createOrderRequestBody)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	testutil.AssertStatus(t, response.Code, http.StatusOK)

	t.Run("get all orders", func(t *testing.T) {
		request := handlers.NewGetAllOrdersRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := []handlers.OrderResponse{
			handlers.NewOrderResponseBody(testdata.PeterCreatedOrder, testdata.PeterCreatedOrderItems, testdata.ChickenShackAddress, testdata.PeterAddress1),
		}
		var got []handlers.OrderResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("get current orders", func(t *testing.T) {
		request := handlers.NewGetCurrentOrdersRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := []handlers.OrderResponse{
			handlers.NewOrderResponseBody(testdata.PeterCreatedOrder, testdata.PeterCreatedOrderItems, testdata.ChickenShackAddress, testdata.PeterAddress1),
		}
		var got []handlers.OrderResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("cancel order", func(t *testing.T) {
		cancelOrderRequestBody := handlers.CancelOrderRequest{ID: 1}
		request := handlers.NewCancelOrderRequest(peterJWT, cancelOrderRequestBody)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := handlers.CancelOrderResponse{Status: true}
		var got handlers.CancelOrderResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("canceled order doesn't show up in current orders", func(t *testing.T) {
		request := handlers.NewGetCurrentOrdersRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := []handlers.OrderResponse{}
		var got []handlers.OrderResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})
}
