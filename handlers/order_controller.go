package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VitoNaychev/bt-order-svc/models"
)

type VerifyJWT func(token string) AuthResponse

type OrderServer struct {
	store     models.OrderStore
	verifyJWT VerifyJWT
}

func (o *OrderServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header["Token"] == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	authResponse := o.verifyJWT(r.Header["Token"][0])
	if authResponse.Status == INVALID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if authResponse.Status == NOT_FOUND {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	order, _ := o.store.GetOrdersByCustomerID(authResponse.ID)
	json.NewEncoder(w).Encode(order)
}
