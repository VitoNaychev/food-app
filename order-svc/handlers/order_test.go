package handlers_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/VitoNaychev/food-app/order-svc/handlers"
	"github.com/VitoNaychev/food-app/order-svc/models"
	"github.com/VitoNaychev/food-app/order-svc/stubs"
	"github.com/VitoNaychev/food-app/order-svc/testdata"
	"github.com/VitoNaychev/food-app/testutil"
	"github.com/VitoNaychev/food-app/testutil/tabletests"
	"github.com/VitoNaychev/food-app/validation"
)

func TestOrderEndpointAuthentication(t *testing.T) {
	orderStore := &stubs.StubOrderStore{}
	orderItemStore := &stubs.StubOrderItemStore{}
	addressStore := &stubs.StubAddressStore{}
	server := handlers.NewOrderServer(orderStore, orderItemStore, addressStore, stubs.StubVerifyJWT)

	invalidJWT := "invalidJWT"
	cases := map[string]*http.Request{
		"get all orders authentication":    handlers.NewGetAllOrdersRequest(invalidJWT),
		"get current order authentication": handlers.NewGetCurrentOrdersRequest(invalidJWT),
		"create order authentication":      handlers.NewCreateOrderRequest(invalidJWT, handlers.CreateOrderRequest{}),
		"cancel order authentication":      handlers.NewCancelOrderRequest(invalidJWT, handlers.CancelOrderRequest{}),
	}

	tabletests.RunAuthenticationTests(t, server, cases)
}

func TestOrderRequestValidation(t *testing.T) {
	orderStore := &stubs.StubOrderStore{}
	orderItemStore := &stubs.StubOrderItemStore{}
	addressStore := &stubs.StubAddressStore{}
	server := handlers.NewOrderServer(orderStore, orderItemStore, addressStore, stubs.StubVerifyJWT)

	peterJWT := strconv.Itoa(testdata.PeterCustomerID)

	cases := map[string]*http.Request{
		"create order": handlers.NewCreateOrderRequest(peterJWT, handlers.CreateOrderRequest{}),
		"cancel order": handlers.NewCancelOrderRequest(peterJWT, handlers.CancelOrderRequest{}),
	}

	tabletests.RunRequestValidationTests(t, server, cases)
}

func TestOrderResponseValidity(t *testing.T) {
	orderStore := &stubs.StubOrderStore{
		Orders: []models.Order{testdata.PeterCreatedOrder, testdata.PeterCompletedOrder},
	}
	orderItemStore := &stubs.StubOrderItemStore{
		OrderItems: concatThreeSlices(testdata.PeterCreatedOrderItems, testdata.PeterCompletedOrderItems, testdata.AliceOrderItems),
	}
	addressStore := &stubs.StubAddressStore{
		Addresses: []models.Address{testdata.ChickenShackAddress, testdata.PeterAddress1, testdata.PeterAddress2},
	}
	server := handlers.NewOrderServer(orderStore, orderItemStore, addressStore, stubs.StubVerifyJWT)

	peterJWT := strconv.Itoa(testdata.PeterCustomerID)
	createOrderRequestBody := handlers.NewCeateOrderRequestBody(testdata.PeterCreatedOrder, testdata.PeterCreatedOrderItems, testdata.ChickenShackAddress, testdata.PeterAddress1)
	cancelOrderRequestBody := handlers.CancelOrderRequest{ID: testdata.PeterCreatedOrder.ID}

	cases := []tabletests.ResponseValidationTestcase{
		{
			Name:    "get all orders",
			Request: handlers.NewGetAllOrdersRequest(peterJWT),
			ValidationFunction: func(r io.Reader) (tabletests.GenericResponse, error) {
				response, err := validation.ValidateBody[[]handlers.OrderResponse](r)
				return response, err
			},
		},
		{
			Name:    "get current orders",
			Request: handlers.NewGetCurrentOrdersRequest(peterJWT),
			ValidationFunction: func(r io.Reader) (tabletests.GenericResponse, error) {
				response, err := validation.ValidateBody[[]handlers.OrderResponse](r)
				return response, err
			},
		},
		{
			Name:    "create order",
			Request: handlers.NewCreateOrderRequest(peterJWT, createOrderRequestBody),
			ValidationFunction: func(r io.Reader) (tabletests.GenericResponse, error) {
				response, err := validation.ValidateBody[handlers.OrderResponse](r)
				return response, err
			},
		},
		{
			Name:    "cancel order",
			Request: handlers.NewCancelOrderRequest(peterJWT, cancelOrderRequestBody),
			ValidationFunction: func(r io.Reader) (tabletests.GenericResponse, error) {
				response, err := validation.ValidateBody[handlers.CancelOrderResponse](r)
				return response, err
			},
		},
	}

	tabletests.RunResponseValidationTests(t, server, cases)
}

func TestCancelOrder(t *testing.T) {
	orderStore := &stubs.StubOrderStore{
		Orders: []models.Order{testdata.PeterCreatedOrder, testdata.PeterCompletedOrder},
	}
	orderItemStore := &stubs.StubOrderItemStore{
		OrderItems: concatThreeSlices(testdata.PeterCreatedOrderItems, testdata.PeterCompletedOrderItems, testdata.AliceOrderItems),
	}
	addressStore := &stubs.StubAddressStore{
		Addresses: []models.Address{testdata.ChickenShackAddress, testdata.PeterAddress1, testdata.PeterAddress2},
	}
	server := handlers.NewOrderServer(orderStore, orderItemStore, addressStore, stubs.StubVerifyJWT)

	t.Run("return Unauthorized on attemp to cancel another user's order", func(t *testing.T) {
		cancelOrderRequestBody := handlers.CancelOrderRequest{ID: 1}
		request := handlers.NewCancelOrderRequest(strconv.Itoa(testdata.AliceCustomerID), cancelOrderRequestBody)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrUnathorizedAction)
	})

	t.Run("returns Status true when order is cancelable", func(t *testing.T) {
		cancelOrderRequestBody := handlers.CancelOrderRequest{ID: 1}
		request := handlers.NewCancelOrderRequest(strconv.Itoa(testdata.PeterCustomerID), cancelOrderRequestBody)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		if orderStore.Orders[0].Status != models.CANCELED {
			t.Errorf("got status %v, want %v", orderStore.Orders[0].Status, models.CANCELED)
		}

		want := handlers.CancelOrderResponse{Status: true}

		var got handlers.CancelOrderResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns Status false when order is noncancelable", func(t *testing.T) {
		cancelOrderRequestBody := handlers.CancelOrderRequest{ID: 2}
		request := handlers.NewCancelOrderRequest(strconv.Itoa(testdata.PeterCustomerID), cancelOrderRequestBody)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		if orderStore.Orders[1].Status == models.CANCELED {
			t.Errorf("status changed to %v when it shouldn't have", orderStore.Orders[0].Status)
		}

		want := handlers.CancelOrderResponse{Status: false}

		var got handlers.CancelOrderResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns Not Found on order that doesn't exist", func(t *testing.T) {
		cancelOrderRequestBody := handlers.CancelOrderRequest{ID: 10}
		request := handlers.NewCancelOrderRequest(strconv.Itoa(testdata.PeterCustomerID), cancelOrderRequestBody)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrOrderNotFound)
	})
}

func TestCreateOrder(t *testing.T) {
	orderStore := &stubs.StubOrderStore{CreatedOrders: []models.Order{}, Orders: nil}
	orderItemStore := &stubs.StubOrderItemStore{CreatedOrderItems: []models.OrderItem{}, OrderItems: nil}
	addressStore := &stubs.StubAddressStore{CreatedAddresses: []models.Address{}, Addresses: nil}
	server := handlers.NewOrderServer(orderStore, orderItemStore, addressStore, stubs.StubVerifyJWT)

	t.Run("creates new order and returns it", func(t *testing.T) {
		createOrderRequestBody := handlers.NewCeateOrderRequestBody(testdata.PeterCreatedOrder, testdata.PeterCreatedOrderItems, testdata.ChickenShackAddress, testdata.PeterAddress1)
		request := handlers.NewCreateOrderRequest(strconv.Itoa(testdata.PeterCustomerID), createOrderRequestBody)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		if len(orderStore.CreatedOrders) != 1 {
			t.Errorf("got %d calls to CreateOrder, want %d", len(orderStore.CreatedOrders), 1)
		}

		if len(addressStore.CreatedAddresses) != 2 {
			t.Errorf("got %d calls to CreateAddress, want %d", len(addressStore.CreatedAddresses), 2)
		}

		if orderStore.CreatedOrders[0].Status != models.APPROVAL_PENDING {
			t.Errorf("got status %v want %v", orderStore.CreatedOrders[0].Status, models.APPROVAL_PENDING)
		}

		wantOrder := testdata.PeterCreatedOrder
		wantOrder.Status = models.APPROVAL_PENDING
		want := handlers.NewOrderResponseBody(wantOrder, testdata.PeterCreatedOrderItems, testdata.ChickenShackAddress, testdata.PeterAddress1)

		var got handlers.OrderResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})
}

func TestGetCurrentOrders(t *testing.T) {
	orderStore := &stubs.StubOrderStore{
		CreatedOrders: nil,
		Orders:        []models.Order{testdata.PeterCreatedOrder, testdata.PeterCompletedOrder, testdata.AliceOrder},
	}
	orderItemStore := &stubs.StubOrderItemStore{
		CreatedOrderItems: nil,
		OrderItems:        concatThreeSlices(testdata.PeterCreatedOrderItems, testdata.PeterCompletedOrderItems, testdata.AliceOrderItems),
	}
	addressStore := &stubs.StubAddressStore{
		CreatedAddresses: nil,
		Addresses:        []models.Address{testdata.ChickenShackAddress, testdata.PeterAddress1, testdata.PeterAddress2, testdata.AliceAddress},
	}
	server := handlers.NewOrderServer(orderStore, orderItemStore, addressStore, stubs.StubVerifyJWT)

	t.Run("returns current orders for customer Peter", func(t *testing.T) {
		request := handlers.NewGetCurrentOrdersRequest(strconv.Itoa(testdata.PeterCustomerID))
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
}

func concatThreeSlices[T any](s1, s2, s3 []T) []T {
	newSlice := make([]T, 0, len(s1)+len(s2)+len(s3))

	newSlice = append(newSlice, append(s1, append(s2, s3...)...)...)
	return newSlice
}

func TestGetOrders(t *testing.T) {
	orderStore := &stubs.StubOrderStore{
		CreatedOrders: nil,
		Orders:        []models.Order{testdata.PeterCreatedOrder, testdata.PeterCompletedOrder, testdata.AliceOrder},
	}
	orderItemStore := &stubs.StubOrderItemStore{
		CreatedOrderItems: nil,
		OrderItems:        concatThreeSlices(testdata.PeterCreatedOrderItems, testdata.PeterCompletedOrderItems, testdata.AliceOrderItems),
	}
	addressStore := &stubs.StubAddressStore{
		CreatedAddresses: nil,
		Addresses:        []models.Address{testdata.ChickenShackAddress, testdata.PeterAddress1, testdata.PeterAddress2, testdata.AliceAddress},
	}
	server := handlers.NewOrderServer(orderStore, orderItemStore, addressStore, stubs.StubVerifyJWT)

	t.Run("returns orders of customer Peter", func(t *testing.T) {
		request := handlers.NewGetAllOrdersRequest(strconv.Itoa(testdata.PeterCustomerID))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := []handlers.OrderResponse{
			handlers.NewOrderResponseBody(testdata.PeterCreatedOrder, testdata.PeterCreatedOrderItems, testdata.ChickenShackAddress, testdata.PeterAddress1),
			handlers.NewOrderResponseBody(testdata.PeterCompletedOrder, testdata.PeterCompletedOrderItems, testdata.ChickenShackAddress, testdata.PeterAddress2),
		}
		var got []handlers.OrderResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns orders of customer Alice", func(t *testing.T) {
		request := handlers.NewGetAllOrdersRequest(strconv.Itoa(testdata.AliceCustomerID))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := []handlers.OrderResponse{
			handlers.NewOrderResponseBody(testdata.AliceOrder, testdata.AliceOrderItems, testdata.ChickenShackAddress, testdata.AliceAddress),
		}

		var got []handlers.OrderResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})
}
