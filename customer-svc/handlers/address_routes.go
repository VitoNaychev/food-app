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
		auth.AuthenticationMiddleware(c.createAddress, c.verifier, c.secretKey)(w, r)
	case http.MethodGet:
		auth.AuthenticationMiddleware(c.getAddress, c.verifier, c.secretKey)(w, r)
	case http.MethodDelete:
		auth.AuthenticationMiddleware(c.deleteAddress, c.verifier, c.secretKey)(w, r)
	case http.MethodPut:
		auth.AuthenticationMiddleware(c.updateAddress, c.verifier, c.secretKey)(w, r)
	}
}
