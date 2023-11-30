package handlers

import (
	"net/http"
	"time"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

type HoursServer struct {
	secretKey       []byte
	expiresAt       time.Duration
	hoursStore      models.HoursStore
	restaurantStore models.RestaurantStore
	verifier        auth.Verifier
}

func NewHoursServer(secretKey []byte, expiresAt time.Duration,
	hoursStore models.HoursStore, restaurantStore models.RestaurantStore) HoursServer {

	hoursServer := HoursServer{
		secretKey:       secretKey,
		expiresAt:       expiresAt,
		hoursStore:      hoursStore,
		restaurantStore: restaurantStore,
		verifier:        NewRestaurantVerifier(restaurantStore),
	}

	return hoursServer
}

func (h *HoursServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		auth.AuthenticationMiddleware(h.getHours, h.verifier, h.secretKey)(w, r)
	case http.MethodPost:
		auth.AuthenticationMiddleware(h.createHours, h.verifier, h.secretKey)(w, r)
	}
}
