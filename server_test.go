package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubOrderStore struct {
	orders []Order
}

func (s *StubOrderStore) GetOrderFromID(id int) (Order, error) {
	for _, order := range s.orders {
		if order.ID == id {
			return order, nil
		}
	}
	return Order{}, fmt.Errorf("No order with ID %d", id)
}

func (s *StubOrderStore) StoreOrder(order Order) (int, error) {
	order.ID = len(s.orders)
	s.orders = append(s.orders, order)

	return order.ID, nil
}

func TestGETOrder(t *testing.T) {
	store := StubOrderStore{[]Order{
		{
			ID:    1,
			Total: 20,
			Items: []string{"taco", "water"},
		},
		{
			ID:    2,
			Total: 40,
			Items: []string{"sushi", "burger", "cola"},
		},
	}}
	server := &OrderServer{&store}

	t.Run("returns 200 on valid order ID", func(t *testing.T) {
		request := newGetOrderRequest(1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("returns 404 on nonexisting order ID", func(t *testing.T) {
		request := newGetOrderRequest(100)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
		// TODO: assert correct error returned
	})

	t.Run("returns 400 on malformed order ID", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/order/adsfasdaf", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("return JSON for order with ID 1", func(t *testing.T) {
		request := newGetOrderRequest(1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getOrderFromResponse(t, response.Body)
		want := Order{
			ID:    1,
			Total: 20,
			Items: []string{"taco", "water"},
		}

		assertStatus(t, response.Code, http.StatusOK)
		assertOrder(t, got, want)
	})

	t.Run("return JSON for order with ID 2", func(t *testing.T) {
		request := newGetOrderRequest(2)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getOrderFromResponse(t, response.Body)
		want := Order{
			ID:    2,
			Total: 40,
			Items: []string{"sushi", "burger", "cola"},
		}

		assertStatus(t, response.Code, http.StatusOK)
		assertOrder(t, got, want)
	})
}

func TestCreateOrder(t *testing.T) {
	store := StubOrderStore{}
	server := &OrderServer{&store}

	t.Run("returns 200 on valid order JSON", func(t *testing.T) {
		order := Order{
			Total: 30,
			Items: []string{"pizza", "pepsi"},
		}

		request := newPostOrderRequest(order)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("returns 400 on invalid order JSON", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte(`{"text":"this is not a valid order JSON"}`))

		request, _ := http.NewRequest("POST", "/order/", buf)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("returns 400 on incomplete order JSON", func(t *testing.T) {
		order := Order{
			Total: 15,
		}

		request := newPostOrderRequest(order)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("persists order and returns order ID", func(t *testing.T) {
		order := Order{
			Total: 15,
			Items: []string{"pizza", "pepsi"},
		}

		request := newPostOrderRequest(order)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertOrderResponse(t, response.Body, CreateOrderResponse{ID: 1})

		order.ID = 1
		assertOrderPersisted(t, &store, order)
	})
}

func getOrderResponseFromResponseBody(body io.Reader) (CreateOrderResponse, error) {
	got := CreateOrderResponse{}
	err := json.NewDecoder(body).Decode(&got)
	if err != nil {
		return CreateOrderResponse{}, fmt.Errorf("unable to unmarshal json, %v", err)
	}

	return got, nil
}

func assertOrderPersisted(t testing.TB, store *StubOrderStore, want Order) {
	t.Helper()

	got, err := store.GetOrderFromID(want.ID)

	if err != nil {
		t.Errorf("server didn't persist order, %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("server didn't persist correct order, got %v want %v", got, want)
	}
}

func assertOrderResponse(t testing.TB, body io.Reader, want CreateOrderResponse) {
	t.Helper()

	got, err := getOrderResponseFromResponseBody(body)
	if err != nil {
		t.Errorf("couldn't get order response from response body,  %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("didn't receive correct order ID, got %v, want %v", got, want)
	}
}

func newPostOrderRequest(order Order) *http.Request {
	buffer := bytes.NewBuffer([]byte{})
	json.NewEncoder(buffer).Encode(order)

	request, _ := http.NewRequest("POST", "/order/", buffer)
	return request
}

func newGetOrderRequest(orderID int) *http.Request {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/order/%d", orderID), nil)
	return req
}

func getOrderFromResponse(t testing.TB, body io.Reader) Order {
	t.Helper()

	var order Order

	err := json.NewDecoder(body).Decode(&order)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into Order, %v", body, err)
	}

	return order
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got status %v, want %v", got, want)
	}
}

func assertOrder(t testing.TB, got, want Order) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
