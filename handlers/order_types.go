package handlers

import (
	"time"

	"github.com/VitoNaychev/bt-order-svc/models"
)

type AuthStatus int

const (
	INVALID AuthStatus = iota
	NOT_FOUND
	OK
)

type AuthResponse struct {
	Status AuthStatus
	ID     int
}

type GetOrderResponse struct {
	ID              int
	CustomerID      int
	RestaurantID    int
	Items           []int
	Total           float64
	DeliveryTime    time.Time
	Status          models.Status
	PickupAddress   models.Address
	DeliveryAddress models.Address
}

func NewGetOrderResponse(order models.Order, pickupAddress, deliveryAddress models.Address) GetOrderResponse {
	return GetOrderResponse{
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

type CreateOrderAddress struct {
	Lat          float64
	Lon          float64
	AddressLine1 string
	AddressLine2 string
	City         string
	Country      string
}
