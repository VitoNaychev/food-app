package handlers

import (
	"net/http"

	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/reqbuilder"
)

func NewChangeDeliveryStateRequest(jwt string, event models.DeliveryEvent) *http.Request {
	stateTransitionDeliveryRequest := StateTransitionDeliveryRequest{
		Event: event,
	}

	request := reqbuilder.NewRequestWithBody(http.MethodPost, "/delivery", stateTransitionDeliveryRequest)
	request.Header.Add("Token", jwt)

	return request
}
