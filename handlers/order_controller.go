package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/VitoNaychev/bt-order-svc/models"
	"github.com/VitoNaychev/errorresponse"
	"github.com/VitoNaychev/validation"
)

var authResponseMap = make(map[string]AuthResponse)

type VerifyJWTFunc func(token string) AuthResponse

func VerifyJWT(token string) (authResponse AuthResponse) {
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:9090/customer/auth/", nil)
	request.Header.Add("Token", token)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("VerifyJWT error: %v", err)
		authResponse.Status = INVALID
		return
	}

	json.NewDecoder(response.Body).Decode(&authResponse)
	return
}

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

func AuthMiddleware(handler func(w http.ResponseWriter, r *http.Request), verifyJWT VerifyJWTFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authResponse := verifyJWT(r.Header["Token"][0])
		if authResponse.Status == INVALID {
			w.WriteHeader(http.StatusUnauthorized)
			writeJSONError(w, http.StatusUnauthorized, ErrInvalidToken)
			return
		}

		if authResponse.Status == NOT_FOUND {
			w.WriteHeader(http.StatusNotFound)
			writeJSONError(w, http.StatusUnauthorized, ErrCustomerNotFound)
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

	err := o.addressStore.CreateAddress(&pickupAddress)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err)
		return
	}
	err = o.addressStore.CreateAddress(&deliveryAddress)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err)
		return
	}

	order.PickupAddress = pickupAddress.ID
	order.DeliveryAddress = deliveryAddress.ID

	order.Status = models.APPROVAL_PENDING
	err = o.orderStore.CreateOrder(&order)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err)
		return
	}

	orderResponse := NewOrderResponseBody(order, pickupAddress, deliveryAddress)

	json.NewEncoder(w).Encode(orderResponse)
}

func (o *OrderServer) getAllOrders(w http.ResponseWriter, r *http.Request) {
	customerJWT := r.Header["Token"][0]
	authResponse := authResponseMap[customerJWT]

	orders, err := o.orderStore.GetOrdersByCustomerID(authResponse.ID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err)
		return
	}

	orderResponseArr := []OrderResponse{}
	for _, order := range orders {
		orderResponseArr = append(orderResponseArr, o.orderToGetOrderResponse(order))
	}

	json.NewEncoder(w).Encode(orderResponseArr)
}

func (o *OrderServer) getCurrentOrders(w http.ResponseWriter, r *http.Request) {
	customerJWT := r.Header["Token"][0]
	authResponse := authResponseMap[customerJWT]

	orders, err := o.orderStore.GetCurrentOrdersByCustomerID(authResponse.ID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err)
	}

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
