package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/VitoNaychev/bt-order-svc/models"
	"github.com/VitoNaychev/errorresponse"
	"github.com/VitoNaychev/validation"
)

var authResponseMap = make(map[string]AuthResponse)

func VerifyJWT(token string) (authResponse AuthResponse, err error) {
	request, err := http.NewRequest(http.MethodPost, "http://customer-svc:8080/customer/auth/", nil)
	if err != nil {
		return
	}
	request.Header.Add("Token", token)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	err = json.NewDecoder(response.Body).Decode(&authResponse)
	return
}

func AuthMiddleware(handler func(w http.ResponseWriter, r *http.Request), verifyJWT VerifyJWTFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authResponse, err := verifyJWT(r.Header["Token"][0])
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}

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

func (o *OrderServer) cancelOrder(w http.ResponseWriter, r *http.Request) {
	var cancelOrderRequest CancelOrderRequest
	json.NewDecoder(r.Body).Decode(&cancelOrderRequest)

	order, err := o.orderStore.GetOrderByID(cancelOrderRequest.ID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, ErrOrderNotFound)
			return
		} else {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
	}
	if order.Status != models.APPROVAL_PENDING && order.Status != models.APPROVED {
		cancelOrderResponse := CancelOrderResponse{Status: false}
		json.NewEncoder(w).Encode(cancelOrderResponse)
		return
	}

	err = o.orderStore.CancelOrder(cancelOrderRequest.ID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err)
		return
	}

	cancelOrderResponse := CancelOrderResponse{Status: true}
	json.NewEncoder(w).Encode(cancelOrderResponse)
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
