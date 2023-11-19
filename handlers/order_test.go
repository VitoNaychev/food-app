package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/VitoNaychev/bt-order-svc/models"
	"github.com/VitoNaychev/bt-order-svc/testdata"
	"github.com/VitoNaychev/bt-order-svc/testutil"
	"github.com/VitoNaychev/validation"
)

func StubVerifyJWT(jwt string) AuthResponse {
	if jwt == "invalidJWT" {
		return AuthResponse{Status: INVALID, ID: 0}
	} else if jwt == "10" {
		return AuthResponse{Status: NOT_FOUND, ID: 0}
	} else {
		id, _ := strconv.Atoi(jwt)
		return AuthResponse{Status: OK, ID: id}
	}
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

func TestAuthMiddleware(t *testing.T) {
	handler := AuthMiddleware(dummyHandler, StubVerifyJWT)

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
	})

	t.Run("returns Not Found on nonexistent customer", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		request.Header.Add("Token", strconv.Itoa(10))
		response := httptest.NewRecorder()

		handler(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
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
	server := NewOrderServer(orderStore, addressStore, StubVerifyJWT)

	invalidJWT := "invalidJWT"
	cases := map[string]*http.Request{
		"get all orders authentication":    NewGetAllOrdersRequest(invalidJWT),
		"get current order authentication": NewGetCurrentOrdersRequest(invalidJWT),
		"create new order authentication":  NewCreateOrderRequest(invalidJWT, CreateOrderRequest{}),
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
	server := NewOrderServer(orderStore, addressStore, StubVerifyJWT)

	peterJWT := strconv.Itoa(testdata.PeterCustomerID)
	createOrderRequestBody := NewCeateOrderRequestBody(testdata.PeterOrder1, testdata.ChickenShackAddress, testdata.PeterAddress1)

	cases := []struct {
		Name               string
		Request            *http.Request
		ValidationFunction func(io.Reader) (GenericResponse, error)
	}{
		{
			"get all orders",
			NewGetAllOrdersRequest(peterJWT),
			func(r io.Reader) (GenericResponse, error) {
				response, err := validation.ValidateBody[[]OrderResponse](r)
				return response, err
			},
		},
		{
			"get current orders",
			NewGetCurrentOrdersRequest(peterJWT),
			func(r io.Reader) (GenericResponse, error) {
				response, err := validation.ValidateBody[[]OrderResponse](r)
				return response, err
			},
		},
		{
			"create order",
			NewCreateOrderRequest(peterJWT, createOrderRequestBody),
			func(r io.Reader) (GenericResponse, error) {
				response, err := validation.ValidateBody[OrderResponse](r)
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

func TestCreateOrder(t *testing.T) {
	orderStore := &testutil.StubOrderStore{CreatedOrders: []models.Order{}, Orders: nil}
	addressStore := &testutil.StubAddressStore{CreatedAddresses: []models.Address{}, Addresses: nil}
	server := NewOrderServer(orderStore, addressStore, StubVerifyJWT)

	t.Run("creates new order and returns it", func(t *testing.T) {
		createOrderRequestBody := NewCeateOrderRequestBody(testdata.PeterOrder1, testdata.ChickenShackAddress, testdata.PeterAddress1)
		request := NewCreateOrderRequest(strconv.Itoa(testdata.PeterCustomerID), createOrderRequestBody)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		if len(orderStore.CreatedOrders) != 1 {
			t.Errorf("got %d calls to CreateOrder, want %d", len(orderStore.CreatedOrders), 1)
		}

		if len(addressStore.CreatedAddresses) != 2 {
			t.Errorf("got %d calls to CreateAddress, want %d", len(addressStore.CreatedAddresses), 2)
		}

		wantOrder := testdata.PeterOrder1
		wantOrder.Status = models.APPROVAL_PENDING
		want := NewOrderResponseBody(wantOrder, testdata.ChickenShackAddress, testdata.PeterAddress1)

		var got OrderResponse
		json.NewDecoder(response.Body).Decode(&got)

		if got.Status != models.APPROVAL_PENDING {
			t.Errorf("got status %v want %v", got.Status, models.APPROVAL_PENDING)
		}

		AssertOrderResponse(t, got, want)
	})
}

func AssertOrderResponse(t testing.TB, got, want OrderResponse) {
	t.Helper()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v want %v", got, want)
	}
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
	server := NewOrderServer(orderStore, addressStore, StubVerifyJWT)

	t.Run("returns current orders for customer with ID 1", func(t *testing.T) {
		request := NewGetCurrentOrdersRequest(strconv.Itoa(testdata.PeterCustomerID))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := []OrderResponse{
			NewOrderResponseBody(testdata.PeterOrder1, testdata.ChickenShackAddress, testdata.PeterAddress1),
		}

		var got []OrderResponse
		json.NewDecoder(response.Body).Decode(&got)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
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
	server := NewOrderServer(orderStore, addressStore, StubVerifyJWT)

	t.Run("returns orders of customer with ID 1", func(t *testing.T) {
		request := NewGetAllOrdersRequest(strconv.Itoa(testdata.PeterCustomerID))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := []OrderResponse{
			NewOrderResponseBody(testdata.PeterOrder1, testdata.ChickenShackAddress, testdata.PeterAddress1),
			NewOrderResponseBody(testdata.PeterOrder2, testdata.ChickenShackAddress, testdata.PeterAddress2),
		}
		var got []OrderResponse
		json.NewDecoder(response.Body).Decode(&got)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("returns orders of customer with ID 2", func(t *testing.T) {
		request := NewGetAllOrdersRequest(strconv.Itoa(testdata.AliceCustomerID))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := []OrderResponse{
			NewOrderResponseBody(testdata.AliceOrder, testdata.ChickenShackAddress, testdata.AliceAddress),
		}

		var got []OrderResponse
		json.NewDecoder(response.Body).Decode(&got)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
