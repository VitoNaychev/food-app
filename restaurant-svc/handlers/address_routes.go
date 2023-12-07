package handlers

import (
	"net/http"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

type AddressServer struct {
	secretKey       []byte
	addressStore    models.AddressStore
	restaurantStore models.RestaurantStore
	verifier        auth.Verifier
}

func NewAddressServer(secretKey []byte, addressStore models.AddressStore, restaurantStore models.RestaurantStore) *AddressServer {
	customerAddressServer := AddressServer{
		secretKey:       secretKey,
		addressStore:    addressStore,
		restaurantStore: restaurantStore,
		verifier:        NewRestaurantVerifier(restaurantStore),
	}

	return &customerAddressServer
}

func (c *AddressServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		auth.AuthenticationMiddleware(c.createAddress, c.verifier, c.secretKey)(w, r)
	case http.MethodGet:
		auth.AuthenticationMiddleware(c.getAddress, c.verifier, c.secretKey)(w, r)
	case http.MethodPut:
		auth.AuthenticationMiddleware(c.updateAddress, c.verifier, c.secretKey)(w, r)
	}
}
