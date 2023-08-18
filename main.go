package main

import (
	"log"
	"net/http"
)

type InMemoryOrderStore struct{}

func (i *InMemoryOrderStore) GetOrderFromID(id int) (Order, error) {
	return Order{ID: 123}, nil
}

func (i *InMemoryOrderStore) StoreOrder(order Order) (int, error) {
	return 1, nil
}

func main() {
	store := &InMemoryOrderStore{}
	server := &OrderServer{store}
	log.Fatal(http.ListenAndServe(":5000", server))
}
