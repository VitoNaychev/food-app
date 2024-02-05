package handlers

import (
	"time"

	"github.com/VitoNaychev/food-app/delivery-svc/models"
)

type DeliveryStateTransitionResponse struct {
	ID    int    `validate:"required,min=1"    json:"id"`
	State string `validate:"required"          json:"state"`
}

func NewDeliveryStateTransitionResponse(delivery models.Delivery) DeliveryStateTransitionResponse {
	stateName, _ := models.StateValueToStateName(delivery.State)

	return DeliveryStateTransitionResponse{
		ID:    delivery.ID,
		State: stateName,
	}
}

type StateTransitionDeliveryRequest struct {
	Event models.DeliveryEvent `validate:"required"          json:"event"`
}

type GetDeliveryResponse struct {
	ID              int                        `validate:"required,min=1"    json:"id"`
	State           string                     `validate:"required"          json:"state"`
	ReadyBy         time.Time                  `validate:"required"          json:"ready_by"`
	PickupAddress   GetDeliveryAddressResponse `validate:"required"          json:"pickup_address"`
	DeliveryAddress GetDeliveryAddressResponse `validate:"required"          json:"delivery_address"`
}

func NewGetDeliveryResponse(delivery models.Delivery, pickupAddress models.Address, deliveryAddress models.Address) GetDeliveryResponse {
	stateName, _ := models.StateValueToStateName(delivery.State)

	pickupAddressResponse := AddressToGetDeliveryAddressResponse(pickupAddress)
	deliveryAddressResponse := AddressToGetDeliveryAddressResponse(deliveryAddress)

	return GetDeliveryResponse{
		ID:              delivery.ID,
		State:           stateName,
		ReadyBy:         delivery.ReadyBy,
		PickupAddress:   pickupAddressResponse,
		DeliveryAddress: deliveryAddressResponse,
	}
}

type GetDeliveryAddressResponse struct {
	Lat          float64 `validate:"latitude,required"  json:"lat"`
	Lon          float64 `validate:"longitude,required" json:"lon"`
	AddressLine1 string  `validate:"required,max=100"   json:"address_line1"`
	AddressLine2 string  `validate:"max=100"            json:"address_line2"`
	City         string  `validate:"required,max=70"    json:"city"`
	Country      string  `validate:"required,max=60"    json:"country"`
}

func AddressToGetDeliveryAddressResponse(address models.Address) GetDeliveryAddressResponse {
	getDeliveryAddressResponse := GetDeliveryAddressResponse{
		Lat:          address.Lat,
		Lon:          address.Lon,
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		City:         address.City,
		Country:      address.Country,
	}

	return getDeliveryAddressResponse
}
