package handlers

import (
	"net/http"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/customer-svc/models"
)

type CustomerAddressServer struct {
	addressStore  models.CustomerAddressStore
	customerStore models.CustomerStore
	secretKey     []byte
	verifier      auth.Verifier
}

func NewCustomerAddressServer(addressStore models.CustomerAddressStore, customerStore models.CustomerStore, secretKey []byte) *CustomerAddressServer {
	customerAddressServer := CustomerAddressServer{
		addressStore:  addressStore,
		customerStore: customerStore,
		secretKey:     secretKey,
		verifier:      NewCustomerVerifier(customerStore),
	}

	return &customerAddressServer
}

func (c *CustomerAddressServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		auth.AuthenticationMW(c.createAddress, c.verifier, c.secretKey)(w, r)
	case http.MethodGet:
		auth.AuthenticationMW(c.getAddress, c.verifier, c.secretKey)(w, r)
	case http.MethodDelete:
		auth.AuthenticationMW(c.deleteAddress, c.verifier, c.secretKey)(w, r)
	case http.MethodPut:
		auth.AuthenticationMW(c.updateAddress, c.verifier, c.secretKey)(w, r)
	}
}
