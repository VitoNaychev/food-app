package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/VitoNaychev/bt-order-svc/handlers"
	"github.com/VitoNaychev/bt-order-svc/models"
	"github.com/VitoNaychev/bt-order-svc/testdata"
	"github.com/VitoNaychev/bt-order-svc/testutil"
)

func TestOrderServerOperations(t *testing.T) {
	connStr := NewDatabaseContainer(t)

	orderStore, err := models.NewPgOrderStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	addressStore, err := models.NewPgAddressStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	server := handlers.NewOrderServer(orderStore, addressStore, testutil.StubVerifyJWT)

	peterJWT := strconv.Itoa(testdata.PeterCustomerID)
	createOrderRequestBody := handlers.NewCeateOrderRequestBody(testdata.PeterOrder1, testdata.ChickenShackAddress, testdata.PeterAddress1)

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
			handlers.NewOrderResponseBody(testdata.PeterOrder1, testdata.ChickenShackAddress, testdata.PeterAddress1),
		}
		var got []handlers.OrderResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertGetOrderResponse(t, got, want)
	})

	t.Run("get current orders", func(t *testing.T) {
		request := handlers.NewGetCurrentOrdersRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := []handlers.OrderResponse{
			handlers.NewOrderResponseBody(testdata.PeterOrder1, testdata.ChickenShackAddress, testdata.PeterAddress1),
		}
		var got []handlers.OrderResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertGetOrderResponse(t, got, want)
	})
}
