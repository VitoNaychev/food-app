package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/VitoNaychev/bt-order-svc/models"
	"github.com/VitoNaychev/bt-order-svc/testdata"
)

type StubOrderStore struct {
	orders []models.Order
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

func StubVerifyJWT(jwt string) AuthResponse {
	if jwt == "invalidJWT" {
		return AuthResponse{INVALID, 0}
	} else if jwt == "10" {
		return AuthResponse{NOT_FOUND, 0}
	} else {
		id, _ := strconv.Atoi(jwt)
		return AuthResponse{OK, id}
	}
}

func TestGetCurrentOrders(t *testing.T) {
	store := &StubOrderStore{[]models.Order{testdata.PeterOrder1, testdata.PeterOrder2, testdata.AliceOrder}}
	server := OrderServer{store, StubVerifyJWT}

	t.Run("returns Unauthorized on missing JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/order/current/", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns Unauthorized on invalid JWT", func(t *testing.T) {
		customerJWT := "invalidJWT"
		request, _ := http.NewRequest(http.MethodGet, "/order/current/", nil)
		request.Header.Add("Token", customerJWT)

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns current orders for customer with ID 1", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/order/current/", nil)
		request.Header.Add("Token", strconv.Itoa(testdata.PeterCustomerID))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		want := []models.Order{testdata.PeterOrder1}

		var got []models.Order
		json.NewDecoder(response.Body).Decode(&got)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestGetOrders(t *testing.T) {
	store := &StubOrderStore{[]models.Order{testdata.PeterOrder1, testdata.PeterOrder2, testdata.AliceOrder}}
	server := OrderServer{store, StubVerifyJWT}

	t.Run("returns Unauthorized on missing JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/order/all/", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns Unauthorized on invalid JWT", func(t *testing.T) {
		customerJWT := "invalidJWT"
		request, _ := http.NewRequest(http.MethodGet, "/order/all/", nil)
		request.Header.Add("Token", customerJWT)

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns orders of customer with ID 1", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/order/all/", nil)
		request.Header.Add("Token", strconv.Itoa(testdata.PeterCustomerID))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		want := []models.Order{testdata.PeterOrder1, testdata.PeterOrder2}

		var got []models.Order
		json.NewDecoder(response.Body).Decode(&got)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("returns orders of customer with ID 2", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/order/all/", nil)
		request.Header.Add("Token", strconv.Itoa(testdata.AliceCustomerID))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		want := []models.Order{testdata.AliceOrder}

		var got []models.Order
		json.NewDecoder(response.Body).Decode(&got)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("returns Not Found on nonexistent customer", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/order/all/", nil)
		request.Header.Add("Token", strconv.Itoa(10))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Fatalf("got status %v want %v", got, want)
	}
}
