package handlers

import (
	"time"

	"github.com/VitoNaychev/bt-order-svc/models"
)

type AuthStatus int

const (
	MISSING_TOKEN AuthStatus = iota
	INVALID
	NOT_FOUND
	OK
)

type AuthResponse struct {
	Status AuthStatus
	ID     int
}

type OrderResponse struct {
	ID              int            `validate:"min=1"`
	CustomerID      int            `validate:"min=1"`
	RestaurantID    int            `validate:"min=1"`
	Items           []int          `validate:"required"`
	Total           float64        `validate:"min=0.01"`
	DeliveryTime    time.Time      `validate:"required"`
	Status          models.Status  `validate:"min=0,max=8"`
	PickupAddress   models.Address `validate:"required"`
	DeliveryAddress models.Address `validate:"required"`
}

func NewOrderResponseBody(order models.Order, pickupAddress, deliveryAddress models.Address) OrderResponse {
	return OrderResponse{
		ID:              order.ID,
		CustomerID:      order.CustomerID,
		RestaurantID:    order.RestaurantID,
		Items:           order.Items,
		Total:           order.Total,
		DeliveryTime:    order.DeliveryTime,
		Status:          order.Status,
		PickupAddress:   pickupAddress,
		DeliveryAddress: deliveryAddress,
	}
}

type CreateOrderRequest struct {
	RestaurantID    int
	Items           []int
	Total           float64
	DeliveryTime    time.Time
	Status          models.Status
	PickupAddress   CreateOrderAddress
	DeliveryAddress CreateOrderAddress
}

func NewCeateOrderRequestBody(order models.Order, pickupAddress models.Address, deliveryAddress models.Address) CreateOrderRequest {
	createPickupAddress := CreateOrderAddress{
		Lat:          pickupAddress.Lat,
		Lon:          pickupAddress.Lon,
		AddressLine1: pickupAddress.AddressLine1,
		AddressLine2: pickupAddress.AddressLine2,
		City:         pickupAddress.City,
		Country:      pickupAddress.Country,
	}

	createDeliveryAddress := CreateOrderAddress{
		Lat:          deliveryAddress.Lat,
		Lon:          deliveryAddress.Lon,
		AddressLine1: deliveryAddress.AddressLine1,
		AddressLine2: deliveryAddress.AddressLine2,
		City:         deliveryAddress.City,
		Country:      deliveryAddress.Country,
	}

	createOrderRequest := CreateOrderRequest{
		RestaurantID:    order.RestaurantID,
		Items:           order.Items,
		Total:           order.Total,
		DeliveryTime:    order.DeliveryTime,
		Status:          order.Status,
		PickupAddress:   createPickupAddress,
		DeliveryAddress: createDeliveryAddress,
	}

	return createOrderRequest
}

func CreateOrderRequestToOrder(createOrderRequest CreateOrderRequest, customerID int) models.Order {
	order := models.Order{
		ID:              0,
		CustomerID:      customerID,
		RestaurantID:    createOrderRequest.RestaurantID,
		Items:           createOrderRequest.Items,
		Total:           createOrderRequest.Total,
		DeliveryTime:    createOrderRequest.DeliveryTime,
		Status:          createOrderRequest.Status,
		PickupAddress:   -1,
		DeliveryAddress: -1,
	}

	return order
}

func GetPickupAddressFromCreateOrderRequest(createOrderRequest CreateOrderRequest) models.Address {
	return CreateOrderAddressToAddress(createOrderRequest.PickupAddress)
}

func GetDeliveryAddressFromCreateOrderRequest(createOrderRequest CreateOrderRequest) models.Address {
	return CreateOrderAddressToAddress(createOrderRequest.DeliveryAddress)
}

type CreateOrderAddress struct {
	Lat          float64
	Lon          float64
	AddressLine1 string
	AddressLine2 string
	City         string
	Country      string
}

func CreateOrderAddressToAddress(createOrderAddress CreateOrderAddress) models.Address {
	address := models.Address{
		ID:           0,
		Lat:          createOrderAddress.Lat,
		Lon:          createOrderAddress.Lon,
		AddressLine1: createOrderAddress.AddressLine1,
		AddressLine2: createOrderAddress.AddressLine2,
		City:         createOrderAddress.City,
		Country:      createOrderAddress.Country,
	}

	return address
}
