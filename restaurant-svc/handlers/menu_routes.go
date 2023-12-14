package handlers

import (
	"net/http"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

type MenuServer struct {
	secretKey       []byte
	menuStore       models.MenuStore
	restaurantStore models.RestaurantStore
	verifier        auth.Verifier
}

func NewMenuServer(secretKey []byte, menuStore models.MenuStore, restaurantStore models.RestaurantStore) *MenuServer {
	return &MenuServer{
		secretKey:       secretKey,
		menuStore:       menuStore,
		restaurantStore: restaurantStore,
		verifier:        NewRestaurantVerifier(restaurantStore),
	}
}

func (m *MenuServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		auth.AuthenticationMW(m.getMenu, m.verifier, m.secretKey)(w, r)
	case http.MethodPost:
		auth.AuthenticationMW(m.createMenuItem, m.verifier, m.secretKey)(w, r)
	case http.MethodPut:
		auth.AuthenticationMW(m.updateMenuItem, m.verifier, m.secretKey)(w, r)
	case http.MethodDelete:
		auth.AuthenticationMW(m.deleteMenuItem, m.verifier, m.secretKey)(w, r)
	}
}
