package handlers

import "github.com/VitoNaychev/food-app/delivery-svc/models"

type StateTransitionDeliveryRequest struct {
	Event models.DeliveryEvent
}
