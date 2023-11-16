package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/VitoNaychev/bt-order-svc/models"
	"github.com/VitoNaychev/bt-order-svc/testdata"
)

type StubAddressStore struct {
	createdAddresses []models.Address
	addresses        []models.Address
}

func (s *StubAddressStore) CreateAddress(address *models.Address) error {
	s.createdAddresses = append(s.createdAddresses, *address)
	address.ID = len(s.createdAddresses)

	return nil
}

func (s *StubAddressStore) GetAddressByID(id int) (models.Address, error) {
	for _, address := range s.addresses {
		if address.ID == id {
			return address, nil
		}
	}
	return models.Address{}, models.ErrNotFound
}

type StubOrderStore struct {
	createdOrders []models.Order
	orders        []models.Order
}

func (s *StubOrderStore) CreateOrder(order *models.Order) error {
	s.createdOrders = append(s.createdOrders, *order)
	order.ID = len(s.createdOrders)

	return nil
}

func (s *StubOrderStore) GetOrdersByCustomerID(customerID int) ([]models.Order, error) {
	var customerOrders []models.Order
	for _, order := range s.orders {
		if order.CustomerID == customerID {
			customerOrders = append(customerOrders, order)
		}
	}
	return customerOrders, nil
}

func (s *StubOrderStore) GetCurrentOrdersByCustomerID(customerID int) ([]models.Order, error) {
	var customerOrders []models.Order
	for _, order := range s.orders {
		if order.CustomerID == customerID && order.Status != models.COMPLETED {
			customerOrders = append(customerOrders, order)
		}
	}
	return customerOrders, nil
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

func stubVerifyJWT(jwt string) AuthResponse {
	if jwt == "invalidJWT" {
		return AuthResponse{INVALID, 0}
	} else if jwt == "10" {
		return AuthResponse{NOT_FOUND, 0}
	} else {
		id, _ := strconv.Atoi(jwt)
		return AuthResponse{OK, id}
	}
}

func TestAuthMiddleware(t *testing.T) {
	handler := AuthMiddleware(dummyHandler, stubVerifyJWT)

	t.Run("returns Unauthorized on missing JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		handler(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns Unauthorized on invalid JWT", func(t *testing.T) {
		customerJWT := "invalidJWT"
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		request.Header.Add("Token", customerJWT)

		response := httptest.NewRecorder()

		handler(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns Not Found on nonexistent customer", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		request.Header.Add("Token", strconv.Itoa(10))
		response := httptest.NewRecorder()

		handler(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns Accepted on authentic customer", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		request.Header.Add("Token", strconv.Itoa(testdata.PeterCustomerID))
		response := httptest.NewRecorder()

		handler(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)
	})
}

func TestOrderEndpointAuthentication(t *testing.T) {
	orderStore := &StubOrderStore{}
	addressStore := &StubAddressStore{}
	server := NewOrderServer(orderStore, addressStore, stubVerifyJWT)

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

			assertStatus(t, response.Code, http.StatusUnauthorized)
		})
	}
}

func TestCreateOrder(t *testing.T) {
	orderStore := &StubOrderStore{[]models.Order{}, nil}
	addressStore := &StubAddressStore{[]models.Address{}, nil}
	server := NewOrderServer(orderStore, addressStore, stubVerifyJWT)

	t.Run("returns Accepted on POST request", func(t *testing.T) {
		createOrderRequestBody := NewCeateOrderRequestBody(testdata.PeterOrder1, testdata.ChickenShackAddress, testdata.PeterAddress1)
		request := NewCreateOrderRequest(strconv.Itoa(testdata.PeterCustomerID), createOrderRequestBody)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		if len(orderStore.createdOrders) != 1 {
			t.Errorf("got %d calls to CreateOrder, want %d", len(orderStore.createdOrders), 1)
		}

		if len(addressStore.createdAddresses) != 2 {
			t.Errorf("got %d calls to CreateAddress, want %d", len(addressStore.createdAddresses), 2)
		}

		want := NewOrderResponseBody(testdata.PeterOrder1, testdata.ChickenShackAddress, testdata.PeterAddress1)
		var got OrderResponse
		json.NewDecoder(response.Body).Decode(&got)

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func NewCeateOrderRequestBody(order models.Order, pickupAddress models.Address, deliveryAddress models.Address) CreateOrderRequest {
	createPickupAddress := CreateOrderAddress{
		Lat:          pickupAddress.Lat,
		Lon:          pickupAddress.Lon,
		AddressLine1: pickupAddress.AddressLine1,
		AddressLine2: pickupAddress.AddressLine2,
		City:         pickupAddress.City,
		Country:      pickupAddress.Country,
	}

	createDeliveryAddress := CreateOrderAddress{
		Lat:          deliveryAddress.Lat,
		Lon:          deliveryAddress.Lon,
		AddressLine1: deliveryAddress.AddressLine1,
		AddressLine2: deliveryAddress.AddressLine2,
		City:         deliveryAddress.City,
		Country:      deliveryAddress.Country,
	}

	createOrderRequest := CreateOrderRequest{
		RestaurantID:    order.RestaurantID,
		Items:           order.Items,
		Total:           order.Total,
		DeliveryTime:    order.DeliveryTime,
		Status:          order.Status,
		PickupAddress:   createPickupAddress,
		DeliveryAddress: createDeliveryAddress,
	}

	return createOrderRequest
}

func NewCreateOrderRequest(jwt string, createOrderRequestBody CreateOrderRequest) *http.Request {
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(createOrderRequestBody)

	request, _ := http.NewRequest(http.MethodPost, "/order/new/", body)
	request.Header.Add("Token", jwt)

	return request
}

func TestGetCurrentOrders(t *testing.T) {
	orderStore := &StubOrderStore{nil, []models.Order{testdata.PeterOrder1, testdata.PeterOrder2, testdata.AliceOrder}}
	addressStore := &StubAddressStore{nil, []models.Address{testdata.ChickenShackAddress, testdata.PeterAddress1, testdata.PeterAddress2, testdata.AliceAddress}}
	server := NewOrderServer(orderStore, addressStore, stubVerifyJWT)

	t.Run("returns current orders for customer with ID 1", func(t *testing.T) {
		request := NewGetCurrentOrdersRequest(strconv.Itoa(testdata.PeterCustomerID))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

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

func NewGetCurrentOrdersRequest(jwt string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/order/current/", nil)
	request.Header.Add("Token", jwt)

	return request
}

func TestGetOrders(t *testing.T) {
	orderStore := &StubOrderStore{nil, []models.Order{testdata.PeterOrder1, testdata.PeterOrder2, testdata.AliceOrder}}
	addressStore := &StubAddressStore{nil, []models.Address{testdata.ChickenShackAddress, testdata.PeterAddress1, testdata.PeterAddress2, testdata.AliceAddress}}
	server := NewOrderServer(orderStore, addressStore, stubVerifyJWT)

	t.Run("returns orders of customer with ID 1", func(t *testing.T) {
		request := NewGetAllOrdersRequest(strconv.Itoa(testdata.PeterCustomerID))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

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

		assertStatus(t, response.Code, http.StatusOK)

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

func assertStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Fatalf("got status %v want %v", got, want)
	}
}

func NewGetAllOrdersRequest(jwt string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/order/all/", nil)
	request.Header.Add("Token", jwt)

	return request
}
