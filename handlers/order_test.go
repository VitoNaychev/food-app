package handlers_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/VitoNaychev/bt-order-svc/handlers"
	"github.com/VitoNaychev/bt-order-svc/models"
	"github.com/VitoNaychev/bt-order-svc/testdata"
	"github.com/VitoNaychev/bt-order-svc/testutil"
	"github.com/VitoNaychev/validation"
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

func TestAuthMiddleware(t *testing.T) {
	handler := handlers.AuthMiddleware(dummyHandler, testutil.StubVerifyJWT)

	t.Run("returns Unauthorized on missing JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		handler(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns Unauthorized on invalid JWT", func(t *testing.T) {
		customerJWT := "invalidJWT"
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		request.Header.Add("Token", customerJWT)

		response := httptest.NewRecorder()

		handler(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrInvalidToken)
	})

	t.Run("returns Not Found on nonexistent customer", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		request.Header.Add("Token", strconv.Itoa(10))
		response := httptest.NewRecorder()

		handler(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrCustomerNotFound)
	})

	t.Run("returns Accepted on authentic customer", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		request.Header.Add("Token", strconv.Itoa(testdata.PeterCustomerID))
		response := httptest.NewRecorder()

		handler(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusAccepted)
	})
}

func TestOrderEndpointAuthentication(t *testing.T) {
	orderStore := &testutil.StubOrderStore{}
	addressStore := &testutil.StubAddressStore{}
	server := handlers.NewOrderServer(orderStore, addressStore, testutil.StubVerifyJWT)

	invalidJWT := "invalidJWT"
	cases := map[string]*http.Request{
		"get all orders authentication":    handlers.NewGetAllOrdersRequest(invalidJWT),
		"get current order authentication": handlers.NewGetCurrentOrdersRequest(invalidJWT),
		"create order authentication":      handlers.NewCreateOrderRequest(invalidJWT, handlers.CreateOrderRequest{}),
		"cancel order authentication":      handlers.NewCancelOrderRequest(invalidJWT, handlers.CancelOrderRequest{}),
	}

	for name, request := range cases {
		t.Run(name, func(t *testing.T) {
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		})
	}
}

type GenericResponse interface{}

func TestOrderResponseValidity(t *testing.T) {
	orderStore := &testutil.StubOrderStore{}
	addressStore := &testutil.StubAddressStore{}
	server := handlers.NewOrderServer(orderStore, addressStore, testutil.StubVerifyJWT)

	peterJWT := strconv.Itoa(testdata.PeterCustomerID)
	createOrderRequestBody := handlers.NewCeateOrderRequestBody(testdata.PeterOrder1, testdata.ChickenShackAddress, testdata.PeterAddress1)

	cases := []struct {
		Name               string
		Request            *http.Request
		ValidationFunction func(io.Reader) (GenericResponse, error)
	}{
		{
			"get all orders",
			handlers.NewGetAllOrdersRequest(peterJWT),
			func(r io.Reader) (GenericResponse, error) {
				response, err := validation.ValidateBody[[]handlers.OrderResponse](r)
				return response, err
			},
		},
		{
			"get current orders",
			handlers.NewGetCurrentOrdersRequest(peterJWT),
			func(r io.Reader) (GenericResponse, error) {
				response, err := validation.ValidateBody[[]handlers.OrderResponse](r)
				return response, err
			},
		},
		{
			"create order",
			handlers.NewCreateOrderRequest(peterJWT, createOrderRequestBody),
			func(r io.Reader) (GenericResponse, error) {
				response, err := validation.ValidateBody[handlers.OrderResponse](r)
				return response, err
			},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			response := httptest.NewRecorder()

			server.ServeHTTP(response, test.Request)

			_, err := test.ValidationFunction(response.Body)
			if err != nil {
				t.Errorf("invalid response body, %v", err)
			}
		})
	}
}

func TestCancelOrder(t *testing.T) {
	orderStore := &testutil.StubOrderStore{
		Orders: []models.Order{testdata.PeterOrder1, testdata.PeterOrder2},
	}
	addressStore := &testutil.StubAddressStore{
		Addresses: []models.Address{testdata.ChickenShackAddress, testdata.PeterAddress1, testdata.PeterAddress2},
	}
	server := handlers.NewOrderServer(orderStore, addressStore, testutil.StubVerifyJWT)

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
	orderStore := &testutil.StubOrderStore{CreatedOrders: []models.Order{}, Orders: nil}
	addressStore := &testutil.StubAddressStore{CreatedAddresses: []models.Address{}, Addresses: nil}
	server := handlers.NewOrderServer(orderStore, addressStore, testutil.StubVerifyJWT)

	t.Run("creates new order and returns it", func(t *testing.T) {
		createOrderRequestBody := handlers.NewCeateOrderRequestBody(testdata.PeterOrder1, testdata.ChickenShackAddress, testdata.PeterAddress1)
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

		wantOrder := testdata.PeterOrder1
		wantOrder.Status = models.APPROVAL_PENDING
		want := handlers.NewOrderResponseBody(wantOrder, testdata.ChickenShackAddress, testdata.PeterAddress1)

		var got handlers.OrderResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})
}

func TestGetCurrentOrders(t *testing.T) {
	orderStore := &testutil.StubOrderStore{
		CreatedOrders: nil,
		Orders:        []models.Order{testdata.PeterOrder1, testdata.PeterOrder2, testdata.AliceOrder},
	}
	addressStore := &testutil.StubAddressStore{
		CreatedAddresses: nil,
		Addresses:        []models.Address{testdata.ChickenShackAddress, testdata.PeterAddress1, testdata.PeterAddress2, testdata.AliceAddress},
	}
	server := handlers.NewOrderServer(orderStore, addressStore, testutil.StubVerifyJWT)

	t.Run("returns current orders for customer Peter", func(t *testing.T) {
		request := handlers.NewGetCurrentOrdersRequest(strconv.Itoa(testdata.PeterCustomerID))
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

func TestGetOrders(t *testing.T) {
	orderStore := &testutil.StubOrderStore{
		CreatedOrders: nil,
		Orders:        []models.Order{testdata.PeterOrder1, testdata.PeterOrder2, testdata.AliceOrder},
	}
	addressStore := &testutil.StubAddressStore{
		CreatedAddresses: nil,
		Addresses:        []models.Address{testdata.ChickenShackAddress, testdata.PeterAddress1, testdata.PeterAddress2, testdata.AliceAddress},
	}
	server := handlers.NewOrderServer(orderStore, addressStore, testutil.StubVerifyJWT)

	t.Run("returns orders of customer Peter", func(t *testing.T) {
		request := handlers.NewGetAllOrdersRequest(strconv.Itoa(testdata.PeterCustomerID))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := []handlers.OrderResponse{
			handlers.NewOrderResponseBody(testdata.PeterOrder1, testdata.ChickenShackAddress, testdata.PeterAddress1),
			handlers.NewOrderResponseBody(testdata.PeterOrder2, testdata.ChickenShackAddress, testdata.PeterAddress2),
		}
		var got []handlers.OrderResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertGetOrderResponse(t, got, want)
	})

	t.Run("returns orders of customer Alice", func(t *testing.T) {
		request := handlers.NewGetAllOrdersRequest(strconv.Itoa(testdata.AliceCustomerID))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := []handlers.OrderResponse{
			handlers.NewOrderResponseBody(testdata.AliceOrder, testdata.ChickenShackAddress, testdata.AliceAddress),
		}

		var got []handlers.OrderResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertGetOrderResponse(t, got, want)
	})
}
