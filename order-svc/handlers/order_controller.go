package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/VitoNaychev/food-app/errorresponse"
	"github.com/VitoNaychev/food-app/msgtypes"
	"github.com/VitoNaychev/food-app/order-svc/models"
	"github.com/VitoNaychev/food-app/validation"
)

var authResponseMap = make(map[string]msgtypes.AuthResponse)

func VerifyJWT(token string) (authResponse msgtypes.AuthResponse, err error) {
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
		if tokenHeader := r.Header.Get("Token"); tokenHeader == "" {
			errorresponse.WriteJSONError(w, http.StatusUnauthorized, ErrMissingToken)
			return
		}

		authResponse, err := verifyJWT(r.Header["Token"][0])
		if err != nil {
			errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
			return
		}

		if authResponse.Status == msgtypes.MISSING_TOKEN {
			errorresponse.WriteJSONError(w, http.StatusUnauthorized, ErrMissingToken)
			return
		}

		if authResponse.Status == msgtypes.INVALID {
			errorresponse.WriteJSONError(w, http.StatusUnauthorized, ErrInvalidToken)
			return
		}

		if authResponse.Status == msgtypes.NOT_FOUND {
			errorresponse.WriteJSONError(w, http.StatusNotFound, ErrCustomerNotFound)
			return
		}

		authResponseMap[r.Header["Token"][0]] = authResponse

		handler(w, r)
	})
}

func (o *OrderServer) cancelOrder(w http.ResponseWriter, r *http.Request) {
	cancelOrderRequest, err := validation.ValidateBody[CancelOrderRequest](r.Body)
	if err != nil {
		errorresponse.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	order, err := o.orderStore.GetOrderByID(cancelOrderRequest.ID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			errorresponse.WriteJSONError(w, http.StatusNotFound, ErrOrderNotFound)
			return
		} else {
			errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
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
		errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	cancelOrderResponse := CancelOrderResponse{Status: true}
	json.NewEncoder(w).Encode(cancelOrderResponse)
}

func (o *OrderServer) createOrder(w http.ResponseWriter, r *http.Request) {
	createOrderRequest, err := validation.ValidateBody[CreateOrderRequest](r.Body)
	if err != nil {
		errorresponse.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	customerJWT := r.Header["Token"][0]
	authResponse := authResponseMap[customerJWT]

	order := CreateOrderRequestToOrder(createOrderRequest, authResponse.ID)
	pickupAddress := GetPickupAddressFromCreateOrderRequest(createOrderRequest)
	deliveryAddress := GetDeliveryAddressFromCreateOrderRequest(createOrderRequest)

	err = o.addressStore.CreateAddress(&pickupAddress)
	if err != nil {
		errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}
	err = o.addressStore.CreateAddress(&deliveryAddress)
	if err != nil {
		errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	order.PickupAddress = pickupAddress.ID
	order.DeliveryAddress = deliveryAddress.ID

	order.Status = models.APPROVAL_PENDING
	err = o.orderStore.CreateOrder(&order)
	if err != nil {
		errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
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
		errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
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
		errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
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
