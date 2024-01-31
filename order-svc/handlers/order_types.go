package handlers

import (
	"time"

	"github.com/VitoNaychev/food-app/events/svcevents"
	"github.com/VitoNaychev/food-app/order-svc/models"
)

func NewOrderCreatedEvent(order models.Order, orderItems []models.OrderItem, pickupAddress, deliveryAddress models.Address) svcevents.OrderCreatedEvent {
	orderCreatedEventPickupAddress := AddressToOrderCreatedEventAddress(pickupAddress)
	orderCreatedEventDeliveryAddress := AddressToOrderCreatedEventAddress(deliveryAddress)

	orderCreatedEventItem := []svcevents.OrderCreatedEventItem{}
	for _, orderItem := range orderItems {
		orderCreatedEventItem = append(orderCreatedEventItem, OrderItemToOrderCreatedEventItem(orderItem))
	}

	orderCreatedEvent := svcevents.OrderCreatedEvent{
		ID:              order.ID,
		RestaurantID:    order.RestaurantID,
		Items:           orderCreatedEventItem,
		Total:           order.Total,
		PickupAddress:   orderCreatedEventPickupAddress,
		DeliveryAddress: orderCreatedEventDeliveryAddress,
	}

	return orderCreatedEvent
}

func OrderItemToOrderCreatedEventItem(orderItem models.OrderItem) svcevents.OrderCreatedEventItem {
	orderCreatedEventItem := svcevents.OrderCreatedEventItem{
		ID:         orderItem.ID,
		MenuItemID: orderItem.MenuItemID,
		Quantity:   orderItem.Quantity,
	}

	return orderCreatedEventItem
}

func AddressToOrderCreatedEventAddress(address models.Address) svcevents.OrderCreatedEventAddress {
	orderCreatedEventAddress := svcevents.OrderCreatedEventAddress{
		ID:           address.ID,
		Lat:          address.Lat,
		Lon:          address.Lon,
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		City:         address.City,
		Country:      address.Country,
	}

	return orderCreatedEventAddress
}

type CancelOrderRequest struct {
	ID int `validate:"min=1" json:"id"`
}

type CancelOrderResponse struct {
	Status bool `json:"status"`
}

type OrderResponse struct {
	ID              int                 `validate:"min=1"       json:"id"`
	CustomerID      int                 `validate:"min=1"       json:"customer_id"`
	RestaurantID    int                 `validate:"min=1"       json:"restaurant_id"`
	Items           []OrderItemResponse `validate:"required"    json:"items"`
	Total           float32             `validate:"min=0.01"    json:"total"`
	DeliveryTime    time.Time           `validate:"required"    json:"delivery_time"`
	Status          models.Status       `validate:"min=0,max=8" json:"status"`
	PickupAddress   AddressResponse     `validate:"required"    json:"pickup_address"`
	DeliveryAddress AddressResponse     `validate:"required"    json:"delivery_address"`
}

func NewOrderResponseBody(order models.Order, orderItems []models.OrderItem, pickupAddress, deliveryAddress models.Address) OrderResponse {
	pickupAddressResponse := AddressToAddressResponse(pickupAddress)
	deliveryAddressResponse := AddressToAddressResponse(deliveryAddress)

	orderItemsResponse := []OrderItemResponse{}
	for _, orderItem := range orderItems {
		orderItemsResponse = append(orderItemsResponse, OrderItemToOrderItemResponse(orderItem))
	}

	return OrderResponse{
		ID:              order.ID,
		CustomerID:      order.CustomerID,
		RestaurantID:    order.RestaurantID,
		Items:           orderItemsResponse,
		Total:           order.Total,
		DeliveryTime:    order.DeliveryTime,
		Status:          order.Status,
		PickupAddress:   pickupAddressResponse,
		DeliveryAddress: deliveryAddressResponse,
	}
}

type OrderItemResponse struct {
	ID         int `validate:"min=1"       json:"id"`
	MenuItemID int `validate:"min=1"       json:"menu_item_id"`
	Quantity   int `validate:"min=1"       json:"quantity"`
}

func OrderItemToOrderItemResponse(orderItem models.OrderItem) OrderItemResponse {
	orderItemResponse := OrderItemResponse{
		ID:         orderItem.ID,
		MenuItemID: orderItem.MenuItemID,
		Quantity:   orderItem.Quantity,
	}

	return orderItemResponse
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
	Items           []CreateOrderItem  `validate:"required" json:"items"`
	Total           float32            `validate:"min=0.01" json:"total"`
	DeliveryTime    time.Time          `validate:"required" json:"delivery_time"`
	PickupAddress   CreateOrderAddress `validate:"required" json:"pickup_address"`
	DeliveryAddress CreateOrderAddress `validate:"required" json:"delivery_address"`
}

func NewCeateOrderRequestBody(order models.Order, orderItems []models.OrderItem, pickupAddress models.Address, deliveryAddress models.Address) CreateOrderRequest {
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

	createOrderItems := []CreateOrderItem{}
	for _, orderItem := range orderItems {
		creatOrderItem := CreateOrderItem{
			MenuItemID: orderItem.MenuItemID,
			Quantity:   orderItem.Quantity,
		}
		createOrderItems = append(createOrderItems, creatOrderItem)
	}

	createOrderRequest := CreateOrderRequest{
		RestaurantID:    order.RestaurantID,
		Items:           createOrderItems,
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
		Total:           createOrderRequest.Total,
		DeliveryTime:    createOrderRequest.DeliveryTime,
		PickupAddress:   -1,
		DeliveryAddress: -1,
	}

	return order
}

func GetOrderItemsFromCreateOrderRequest(createOrderRequest CreateOrderRequest) []models.OrderItem {
	orderItems := []models.OrderItem{}

	for _, createOrderItem := range createOrderRequest.Items {
		orderItems = append(orderItems, CreateOrderItemToOrderItem(createOrderItem))
	}

	return orderItems
}

type CreateOrderItem struct {
	MenuItemID int `validate:"min=1"       json:"menu_item_id"`
	Quantity   int `validate:"min=1"       json:"quantity"`
}

func CreateOrderItemToOrderItem(createOrderItem CreateOrderItem) models.OrderItem {
	orderItem := models.OrderItem{
		MenuItemID: createOrderItem.MenuItemID,
		Quantity:   createOrderItem.Quantity,
	}

	return orderItem
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
