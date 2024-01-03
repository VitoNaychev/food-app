package handlers

import (
	"net/http"
	"time"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

type RestaurantServer struct {
	secretKey []byte
	expiresAt time.Duration
	store     models.RestaurantStore
	verifier  auth.Verifier
	publisher events.EventPublisher
	http.Handler
}

func NewRestaurantServer(secretKey []byte, expiresAt time.Duration, store models.RestaurantStore, publisher events.EventPublisher) *RestaurantServer {
	s := RestaurantServer{
		secretKey: secretKey,
		expiresAt: expiresAt,
		store:     store,
		verifier:  NewRestaurantVerifier(store),
		publisher: publisher,
	}

	router := http.NewServeMux()
	router.HandleFunc("/restaurant/", s.RestaurantHandler)
	router.HandleFunc("/restaurant/login/", s.LoginHandler)

	s.Handler = router

	return &s
}

func (s *RestaurantServer) RestaurantHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.createRestaurant(w, r)
	case http.MethodGet:
		auth.AuthenticationMW(s.getRestaurant, s.verifier, s.secretKey)(w, r)
	case http.MethodPut:
		auth.AuthenticationMW(s.updateRestaurant, s.verifier, s.secretKey)(w, r)
	case http.MethodDelete:
		auth.AuthenticationMW(s.deleteRestaurant, s.verifier, s.secretKey)(w, r)
	}
}
