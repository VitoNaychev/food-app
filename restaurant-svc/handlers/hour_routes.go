package handlers

import (
	"net/http"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

type HoursServer struct {
	secretKey       []byte
	hoursStore      models.HoursStore
	restaurantStore models.RestaurantStore
	verifier        auth.Verifier
}

func NewHoursServer(secretKey []byte,
	hoursStore models.HoursStore, restaurantStore models.RestaurantStore) *HoursServer {

	hoursServer := &HoursServer{
		secretKey:       secretKey,
		hoursStore:      hoursStore,
		restaurantStore: restaurantStore,
		verifier:        NewRestaurantVerifier(restaurantStore),
	}

	return hoursServer
}

func (h *HoursServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		auth.AuthenticationMW(h.getHours, h.verifier, h.secretKey)(w, r)
	case http.MethodPost:
		auth.AuthenticationMW(h.createHours, h.verifier, h.secretKey)(w, r)
	case http.MethodPut:
		auth.AuthenticationMW(h.updateHours, h.verifier, h.secretKey)(w, r)
	}
}
