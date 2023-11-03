package address

import (
	"net/http"

	"github.com/VitoNaychev/bt-customer-svc/handlers/auth"
	"github.com/VitoNaychev/bt-customer-svc/models"
)

type CustomerAddressServer struct {
	addressStore  models.CustomerAddressStore
	customerStore models.CustomerStore
	secretKey     []byte
}

func NewCustomerAddressServer(addressStore models.CustomerAddressStore, customerStore models.CustomerStore, secretKey []byte) *CustomerAddressServer {
	customerAddressServer := CustomerAddressServer{
		addressStore:  addressStore,
		customerStore: customerStore,
		secretKey:     secretKey,
	}

	return &customerAddressServer
}

func (c *CustomerAddressServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		auth.AuthenticationMiddleware(c.createAddress, c.secretKey)(w, r)
	case http.MethodGet:
		auth.AuthenticationMiddleware(c.getAddress, c.secretKey)(w, r)
	case http.MethodDelete:
		auth.AuthenticationMiddleware(c.deleteAddress, c.secretKey)(w, r)
	case http.MethodPut:
		auth.AuthenticationMiddleware(c.updateAddress, c.secretKey)(w, r)
	}
}
