package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VitoNaychev/bt-order-svc/models"
	"github.com/VitoNaychev/validation"
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
	router.Handle("/order/new/", AuthMiddleware(server.createOrder, verifyJWT))

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

func (o *OrderServer) createOrder(w http.ResponseWriter, r *http.Request) {
	createOrderRequest, _ := validation.ValidateBody[CreateOrderRequest](r.Body)

	customerJWT := r.Header["Token"][0]
	authResponse := authResponseMap[customerJWT]

	order := CreateOrderRequestToOrder(createOrderRequest, authResponse.ID)
	pickupAddress := GetPickupAddressFromCreateOrderRequest(createOrderRequest)
	deliveryAddress := GetDeliveryAddressFromCreateOrderRequest(createOrderRequest)

	_ = o.addressStore.CreateAddress(&pickupAddress)
	_ = o.addressStore.CreateAddress(&deliveryAddress)

	order.PickupAddress = pickupAddress.ID
	order.DeliveryAddress = deliveryAddress.ID

	order.Status = models.APPROVAL_PENDING
	_ = o.orderStore.CreateOrder(&order)

	orderResponse := NewOrderResponseBody(order, pickupAddress, deliveryAddress)

	json.NewEncoder(w).Encode(orderResponse)
}

func (o *OrderServer) getAllOrders(w http.ResponseWriter, r *http.Request) {
	customerJWT := r.Header["Token"][0]
	authResponse := authResponseMap[customerJWT]

	orders, _ := o.orderStore.GetOrdersByCustomerID(authResponse.ID)

	orderResponseArr := []OrderResponse{}
	for _, order := range orders {
		orderResponseArr = append(orderResponseArr, o.orderToGetOrderResponse(order))
	}

	json.NewEncoder(w).Encode(orderResponseArr)
}

func (o *OrderServer) getCurrentOrders(w http.ResponseWriter, r *http.Request) {
	customerJWT := r.Header["Token"][0]
	authResponse := authResponseMap[customerJWT]

	orders, _ := o.orderStore.GetCurrentOrdersByCustomerID(authResponse.ID)

	orderResponseArr := []OrderResponse{}
	for _, order := range orders {
		orderResponseArr = append(orderResponseArr, o.orderToGetOrderResponse(order))
	}

	json.NewEncoder(w).Encode(orderResponseArr)
}

func (o *OrderServer) orderToGetOrderResponse(order models.Order) OrderResponse {
	pickupAddress, _ := o.addressStore.GetAddressByID(order.PickupAddress)
	deliveryAddress, _ := o.addressStore.GetAddressByID(order.DeliveryAddress)

	getOrderResponse := NewOrderResponseBody(order, pickupAddress, deliveryAddress)
	return getOrderResponse
}

func (o *OrderServer) getAllOrders(w http.ResponseWriter, r *http.Request, authResponse AuthResponse) {
	orders, _ := o.store.GetOrdersByCustomerID(authResponse.ID)
	json.NewEncoder(w).Encode(orders)
}

func (o *OrderServer) getCurrentOrders(w http.ResponseWriter, r *http.Request, authResponse AuthResponse) {
	orders, _ := o.store.GetCurrentOrdersByCustomerID(authResponse.ID)
	json.NewEncoder(w).Encode(orders)
}
