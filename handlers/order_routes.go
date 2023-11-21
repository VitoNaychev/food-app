package handlers

import (
	"net/http"

	"github.com/VitoNaychev/bt-order-svc/models"
)

type OrderServer struct {
	orderStore   models.OrderStore
	addressStore models.AddressStore
	verifyJWT    VerifyJWTFunc
	http.Handler
}

func NewOrderServer(orderStore models.OrderStore, addressStore models.AddressStore, verifyJWT VerifyJWTFunc) OrderServer {
	server := OrderServer{
		orderStore:   orderStore,
		addressStore: addressStore,
		verifyJWT:    verifyJWT,
	}

	router := http.NewServeMux()

	router.Handle("/order/all/", AuthMiddleware(server.getAllOrders, verifyJWT))
	router.Handle("/order/current/", AuthMiddleware(server.getCurrentOrders, verifyJWT))
	router.Handle("/order/new/", AuthMiddleware(server.createOrder, verifyJWT))

	server.Handler = router

	return server
}
