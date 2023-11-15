package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VitoNaychev/bt-order-svc/models"
)

var authResponseMap = make(map[string]AuthResponse)

type VerifyJWT func(token string) AuthResponse

type OrderServer struct {
	orderStore   models.OrderStore
	addressStore models.AddressStore
	verifyJWT    VerifyJWT
	http.Handler
}

func NewOrderServer(orderStore models.OrderStore, addressStore models.AddressStore, verifyJWT VerifyJWT) OrderServer {
	server := OrderServer{
		orderStore:   orderStore,
		addressStore: addressStore,
		verifyJWT:    verifyJWT,
	}

	router := http.NewServeMux()

	router.Handle("/order/all/", AuthMiddleware(server.getAllOrders, verifyJWT))
	router.Handle("/order/current/", AuthMiddleware(server.getCurrentOrders, verifyJWT))

	server.Handler = router

	return server
}

func AuthMiddleware(handler func(w http.ResponseWriter, r *http.Request), verifyJWT VerifyJWT) http.HandlerFunc {
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

	orders, _ := o.orderStore.GetOrdersByCustomerID(authResponse.ID)

	orderResponseArr := []GetOrderResponse{}
	for _, order := range orders {
		orderResponseArr = append(orderResponseArr, o.orderToGetOrderResponse(order))
	}

	json.NewEncoder(w).Encode(orderResponseArr)
}

func (o *OrderServer) getCurrentOrders(w http.ResponseWriter, r *http.Request) {
	customerJWT := r.Header["Token"][0]
	authResponse := authResponseMap[customerJWT]

	orders, _ := o.orderStore.GetCurrentOrdersByCustomerID(authResponse.ID)

	orderResponseArr := []GetOrderResponse{}
	for _, order := range orders {
		orderResponseArr = append(orderResponseArr, o.orderToGetOrderResponse(order))
	}

	json.NewEncoder(w).Encode(orderResponseArr)
}

func (o *OrderServer) orderToGetOrderResponse(order models.Order) GetOrderResponse {
	pickupAddress, _ := o.addressStore.GetAddressByID(order.PickupAddress)
	deliveryAddress, _ := o.addressStore.GetAddressByID(order.DeliveryAddress)

	getOrderResponse := NewGetOrderResponse(order, pickupAddress, deliveryAddress)
	return getOrderResponse
}
