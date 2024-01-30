package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/food-app/httperrors"
	"github.com/VitoNaychev/food-app/msgtypes"
	"github.com/VitoNaychev/food-app/order-svc/models"
	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/VitoNaychev/food-app/validation"
)

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

func (o *OrderServer) cancelOrder(w http.ResponseWriter, r *http.Request) {
	cancelOrderRequest, err := validation.ValidateBody[CancelOrderRequest](r.Body)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	order, err := o.orderStore.GetOrderByID(cancelOrderRequest.ID)
	if err != nil {
		if errors.Is(err, storeerrors.ErrNotFound) {
			httperrors.WriteJSONError(w, http.StatusNotFound, ErrOrderNotFound)
			return
		} else {
			httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
			return
		}
	}

	customerID, _ := strconv.Atoi(r.Header["Subject"][0])
	if order.CustomerID != customerID {
		httperrors.HandleUnauthorized(w, ErrUnathorizedAction)
		return
	}

	if order.Status != models.APPROVAL_PENDING && order.Status != models.APPROVED {
		cancelOrderResponse := CancelOrderResponse{Status: false}
		json.NewEncoder(w).Encode(cancelOrderResponse)
		return
	}

	err = o.orderStore.CancelOrder(cancelOrderRequest.ID)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	cancelOrderResponse := CancelOrderResponse{Status: true}
	json.NewEncoder(w).Encode(cancelOrderResponse)
}

func (o *OrderServer) createOrder(w http.ResponseWriter, r *http.Request) {
	createOrderRequest, err := validation.ValidateBody[CreateOrderRequest](r.Body)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	customerID, _ := strconv.Atoi(r.Header["Subject"][0])

	order := CreateOrderRequestToOrder(createOrderRequest, customerID)
	pickupAddress := GetPickupAddressFromCreateOrderRequest(createOrderRequest)
	deliveryAddress := GetDeliveryAddressFromCreateOrderRequest(createOrderRequest)

	err = o.addressStore.CreateAddress(&pickupAddress)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}
	err = o.addressStore.CreateAddress(&deliveryAddress)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	order.PickupAddress = pickupAddress.ID
	order.DeliveryAddress = deliveryAddress.ID

	order.Status = models.APPROVAL_PENDING
	err = o.orderStore.CreateOrder(&order)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	orderItems := GetOrderItemsFromCreateOrderRequest(createOrderRequest)
	for i := range orderItems {
		orderItems[i].OrderID = order.ID

		err = o.orderItemStore.CreateOrderItem(&orderItems[i])
		if err != nil {
			httperrors.HandleInternalServerError(w, err)
			return
		}
	}

	orderResponse := NewOrderResponseBody(order, orderItems, pickupAddress, deliveryAddress)

	json.NewEncoder(w).Encode(orderResponse)
}

func (o *OrderServer) getAllOrders(w http.ResponseWriter, r *http.Request) {
	customerID, _ := strconv.Atoi(r.Header["Subject"][0])

	orders, err := o.orderStore.GetOrdersByCustomerID(customerID)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	orderResponseArr := []OrderResponse{}
	for _, order := range orders {
		orderResponseArr = append(orderResponseArr, o.orderToGetOrderResponse(order))
	}

	json.NewEncoder(w).Encode(orderResponseArr)
}

func (o *OrderServer) getCurrentOrders(w http.ResponseWriter, r *http.Request) {
	customerID, _ := strconv.Atoi(r.Header["Subject"][0])

	orders, err := o.orderStore.GetCurrentOrdersByCustomerID(customerID)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
	}

	orderResponseArr := []OrderResponse{}
	for _, order := range orders {
		orderResponseArr = append(orderResponseArr, o.orderToGetOrderResponse(order))
	}

	json.NewEncoder(w).Encode(orderResponseArr)
}

func (o *OrderServer) orderToGetOrderResponse(order models.Order) OrderResponse {
	orderItems, _ := o.orderItemStore.GetOrderItemsByOrderID(order.ID)

	pickupAddress, _ := o.addressStore.GetAddressByID(order.PickupAddress)
	deliveryAddress, _ := o.addressStore.GetAddressByID(order.DeliveryAddress)

	getOrderResponse := NewOrderResponseBody(order, orderItems, pickupAddress, deliveryAddress)
	return getOrderResponse
}
