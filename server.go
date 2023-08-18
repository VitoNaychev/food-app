package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Order struct {
	ID    int
	Total int
	Items []string
}

type CreateOrderResponse struct {
	ID int
}

type OrderStore interface {
	GetOrderFromID(id int) (Order, error)
	StoreOrder(order Order) (int, error)
}

type OrderServer struct {
	store OrderStore
}

func (o *OrderServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		o.getOrder(w, r)
	case http.MethodPost:
		o.createOrder(w, r)
	}

}

func (o *OrderServer) getOrder(w http.ResponseWriter, r *http.Request) {
	orderID, err := getOrderIDFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	order, err := o.store.GetOrderFromID(orderID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (o *OrderServer) createOrder(w http.ResponseWriter, r *http.Request) {
	order, err := getOrderFromRequestBody(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := o.store.StoreOrder(*order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CreateOrderResponse{ID: id})
}

func getOrderFromRequestBody(body io.Reader) (*Order, error) {
	d := json.NewDecoder(body)
	d.DisallowUnknownFields()

	order := Order{}
	err := d.Decode(&order)

	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal order JSON, %v", err)
	}

	if !isOrderValid(order) {
		return nil, fmt.Errorf("some fields of order JSON are empty, cannot persist")
	}

	return &order, nil
}

func isOrderValid(order Order) bool {
	if order.Total == 0 {
		return false
	}

	if order.Items == nil {
		return false
	}

	return true
}

func getOrderIDFromRequest(r *http.Request) (int, error) {
	stringID := strings.TrimPrefix(r.URL.Path, "/order/")
	return strconv.Atoi(stringID)
}
