package address

import (
	"net/http"

	"github.com/VitoNaychev/bt-customer-svc/handlers/auth"
	"github.com/VitoNaychev/bt-customer-svc/models/address_store"
	"github.com/VitoNaychev/bt-customer-svc/models/customer_store"
)

type CustomerAddressServer struct {
	addressStore  address_store.CustomerAddressStore
	customerStore customer_store.CustomerStore
	secretKey     []byte
}

func (c *CustomerAddressServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		auth.AuthenticationMiddleware(c.StoreAddressHandler, c.secretKey)(w, r)
	case http.MethodGet:
		auth.AuthenticationMiddleware(c.GetAddressHandler, c.secretKey)(w, r)
	case http.MethodDelete:
		auth.AuthenticationMiddleware(c.DeleteAddressHandler, c.secretKey)(w, r)
	case http.MethodPut:
		auth.AuthenticationMiddleware(c.UpdateAddress, c.secretKey)(w, r)
	}
}
