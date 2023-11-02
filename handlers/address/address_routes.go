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
