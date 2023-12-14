package handlers

import (
	"net/http"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/order-svc/models"
)

type OrderServer struct {
	orderStore   models.OrderStore
	addressStore models.AddressStore
	verifyJWT    auth.VerifyJWTFunc
	http.Handler
}

func NewOrderServer(orderStore models.OrderStore, addressStore models.AddressStore, verifyJWT auth.VerifyJWTFunc) OrderServer {
	server := OrderServer{
		orderStore:   orderStore,
		addressStore: addressStore,
		verifyJWT:    verifyJWT,
	}

	router := http.NewServeMux()

	router.Handle("/order/all/", auth.RemoteAuthenticationMW(server.getAllOrders, verifyJWT))
	router.Handle("/order/current/", auth.RemoteAuthenticationMW(server.getCurrentOrders, verifyJWT))
	router.Handle("/order/new/", auth.RemoteAuthenticationMW(server.createOrder, verifyJWT))
	router.Handle("/order/cancel/", auth.RemoteAuthenticationMW(server.cancelOrder, verifyJWT))

	server.Handler = router

	return server
}
