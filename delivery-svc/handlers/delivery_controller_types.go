package handlers

import (
	"time"

	"github.com/VitoNaychev/food-app/delivery-svc/models"
)

type DeliveryStateTransitionResponse struct {
	ID    int
	State string
}

func NewDeliveryStateTransitionResponse(delivery models.Delivery) DeliveryStateTransitionResponse {
	stateName, _ := models.StateValueToStateName(delivery.State)

	return DeliveryStateTransitionResponse{
		ID:    delivery.ID,
		State: stateName,
	}
}

type StateTransitionDeliveryRequest struct {
	Event models.DeliveryEvent
}

type GetDeliveryResponse struct {
	ID              int
	State           string
	ReadyBy         time.Time
	PickupAddress   models.Address
	DeliveryAddress models.Address
}

func NewGetDeliveryResponse(delivery models.Delivery, pickupAddress models.Address, deliveryAddress models.Address) GetDeliveryResponse {
	stateName, _ := models.StateValueToStateName(delivery.State)

	return GetDeliveryResponse{
		ID:              delivery.ID,
		State:           stateName,
		ReadyBy:         delivery.ReadyBy,
		PickupAddress:   pickupAddress,
		DeliveryAddress: deliveryAddress,
	}
}
