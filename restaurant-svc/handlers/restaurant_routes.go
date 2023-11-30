package handlers

import (
	"net/http"
	"time"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

type RestaurantServer struct {
	secretKey []byte
	expiresAt time.Duration
	store     models.RestaurantStore
}

func NewRestaurantserver(secretKey []byte, expiresAt time.Duration, store models.RestaurantStore) RestaurantServer {
	restaurantServer := RestaurantServer{
		secretKey: secretKey,
		expiresAt: expiresAt,
		store:     store,
	}

	return restaurantServer
}

func (s *RestaurantServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.createRestaurant(w, r)
	case http.MethodGet:
		auth.AuthenticationMiddleware(s.getRestaurant, s.secretKey)(w, r)
	case http.MethodPut:
		auth.AuthenticationMiddleware(s.updateRestaurant, s.secretKey)(w, r)
	}
}
