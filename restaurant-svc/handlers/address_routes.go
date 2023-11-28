package handlers

import (
	"net/http"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

type AddressServer struct {
	addressStore    models.AddressStore
	restaurantStore models.RestaurantStore
	secretKey       []byte
}

func NewAddressServer(addressStore models.AddressStore, restaurantStore models.RestaurantStore, secretKey []byte) *AddressServer {
	customerAddressServer := AddressServer{
		addressStore:    addressStore,
		restaurantStore: restaurantStore,
		secretKey:       secretKey,
	}

	return &customerAddressServer
}

func (c *AddressServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		auth.AuthenticationMiddleware(c.createAddress, c.secretKey)(w, r)
	case http.MethodGet:
		auth.AuthenticationMiddleware(c.getAddress, c.secretKey)(w, r)
	case http.MethodPut:
		auth.AuthenticationMiddleware(c.updateAddress, c.secretKey)(w, r)
	}
}
