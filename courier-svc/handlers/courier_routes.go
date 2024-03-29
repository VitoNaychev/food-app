package handlers

import (
	"net/http"
	"time"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/courier-svc/models"
	"github.com/VitoNaychev/food-app/events"
)

type CourierServer struct {
	secretKey []byte
	expiresAt time.Duration
	store     models.CourierStore
	publisher events.EventPublisher

	verifier auth.Verifier
	http.Handler
}

func NewCourierServer(secretKey []byte, expiresAt time.Duration, store models.CourierStore, publisher events.EventPublisher) *CourierServer {
	s := CourierServer{
		secretKey: secretKey,
		expiresAt: expiresAt,
		store:     store,
		publisher: publisher,

		verifier: NewCourierVerifier(store),
	}

	router := http.NewServeMux()
	router.HandleFunc("/courier/", s.CourierHandler)
	router.HandleFunc("/courier/login/", s.LoginHandler)

	s.Handler = router

	return &s
}

func (s *CourierServer) CourierHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.createCourier(w, r)
	case http.MethodGet:
		auth.AuthenticationMW(s.getCourier, s.verifier, s.secretKey)(w, r)
	case http.MethodPut:
		auth.AuthenticationMW(s.updateCourier, s.verifier, s.secretKey)(w, r)
	case http.MethodDelete:
		auth.AuthenticationMW(s.deleteCourier, s.verifier, s.secretKey)(w, r)
	}
}
