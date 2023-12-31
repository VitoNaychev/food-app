package handlers

import (
	"time"

	"github.com/VitoNaychev/food-app/order-svc/models"
)

type CancelOrderRequest struct {
	ID int `validate:"min=1" json:"id"`
}

type CancelOrderResponse struct {
	Status bool `json:"status"`
}

type OrderResponse struct {
	ID              int             `validate:"min=1"       json:"id"`
	CustomerID      int             `validate:"min=1"       json:"customer_id"`
	RestaurantID    int             `validate:"min=1"       json:"restaurant_id"`
	Items           []int           `validate:"required"    json:"items"`
	Total           float64         `validate:"min=0.01"    json:"total"`
	DeliveryTime    time.Time       `validate:"required"    json:"delivery_time"`
	Status          models.Status   `validate:"min=0,max=8" json:"status"`
	PickupAddress   AddressResponse `validate:"required"    json:"pickup_address"`
	DeliveryAddress AddressResponse `validate:"required"    json:"delivery_address"`
}

type AddressResponse struct {
	Id           int     `validate:"min=1"               json:"id"`
	Lat          float64 `validate:"latitude,required"   json:"lat"`
	Lon          float64 `validate:"longitude,required"  json:"lon"`
	AddressLine1 string  `validate:"required,max=100"    json:"address_line1"`
	AddressLine2 string  `validate:"max=100"             json:"address_line2"`
	City         string  `validate:"required,max=70"     json:"city"`
	Country      string  `validate:"required,max=60"     json:"country"`
}

func NewOrderResponseBody(order models.Order, pickupAddress, deliveryAddress models.Address) OrderResponse {
	pickupAddressResponse := AddressToAddressResponse(pickupAddress)
	deliveryAddressResponse := AddressToAddressResponse(deliveryAddress)

	return OrderResponse{
		ID:              order.ID,
		CustomerID:      order.CustomerID,
		RestaurantID:    order.RestaurantID,
		Items:           order.Items,
		Total:           order.Total,
		DeliveryTime:    order.DeliveryTime,
		Status:          order.Status,
		PickupAddress:   pickupAddressResponse,
		DeliveryAddress: deliveryAddressResponse,
	}
}

func AddressToAddressResponse(address models.Address) AddressResponse {
	return AddressResponse{
		Id:           address.ID,
		Lat:          address.Lat,
		Lon:          address.Lon,
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		City:         address.City,
		Country:      address.Country,
	}
}

type CreateOrderRequest struct {
	RestaurantID    int                `validate:"min=1"    json:"restaurant_id"`
	Items           []int              `validate:"required" json:"items"`
	Total           float64            `validate:"min=0.01" json:"total"`
	DeliveryTime    time.Time          `validate:"required" json:"delivery_time"`
	PickupAddress   CreateOrderAddress `validate:"required" json:"pickup_address"`
	DeliveryAddress CreateOrderAddress `validate:"required" json:"delivery_address"`
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
	Lat          float64 `validate:"required,latitude" json:"lat"`
	Lon          float64 `validate:"required,latitude" json:"lon"`
	AddressLine1 string  `validate:"required,max=100"  json:"address_line1"`
	AddressLine2 string  `validate:"max=100"           json:"address_line2"`
	City         string  `validate:"required,max=70"   json:"city"`
	Country      string  `validate:"required,max=60"   json:"country"`
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
