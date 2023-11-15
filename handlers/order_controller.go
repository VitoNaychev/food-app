package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VitoNaychev/bt-order-svc/models"
)

var authResponseMap = make(map[string]AuthResponse)

type VerifyJWT func(token string) AuthResponse

type OrderServer struct {
	store     models.OrderStore
	verifyJWT VerifyJWT
	http.Handler
}

func NewOrderServer(store models.OrderStore, verifyJWT VerifyJWT) OrderServer {
	server := OrderServer{
		store:     store,
		verifyJWT: verifyJWT,
	}

	router := http.NewServeMux()

	router.Handle("/order/all/", AuthenticationMiddleware(server.getAllOrders, verifyJWT))
	router.Handle("/order/current/", AuthenticationMiddleware(server.getCurrentOrders, verifyJWT))

	server.Handler = router

	return server
}

func AuthenticationMiddleware(handler func(w http.ResponseWriter, r *http.Request), verifyJWT VerifyJWT) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authResponse := verifyJWT(r.Header["Token"][0])
		if authResponse.Status == INVALID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if authResponse.Status == NOT_FOUND {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		authResponseMap[r.Header["Token"][0]] = authResponse

		handler(w, r)
	})
}

func (o *OrderServer) getAllOrders(w http.ResponseWriter, r *http.Request) {
	customerJWT := r.Header["Token"][0]
	authResponse := authResponseMap[customerJWT]

	orders, _ := o.store.GetOrdersByCustomerID(authResponse.ID)
	json.NewEncoder(w).Encode(orders)
}

func (o *OrderServer) getCurrentOrders(w http.ResponseWriter, r *http.Request) {
	customerJWT := r.Header["Token"][0]
	authResponse := authResponseMap[customerJWT]

	orders, _ := o.store.GetCurrentOrdersByCustomerID(authResponse.ID)
	json.NewEncoder(w).Encode(orders)
}
