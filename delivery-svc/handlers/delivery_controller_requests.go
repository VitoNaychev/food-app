package handlers

import (
	"net/http"

	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/reqbuilder"
)

func NewGetActiveDeliveryRequest(jwt string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/delivery/", nil)
	request.Header.Add("Token", jwt)

	return request
}

func NewChangeDeliveryStateRequest(jwt string, event models.DeliveryEvent) *http.Request {
	stateTransitionDeliveryRequest := StateTransitionDeliveryRequest{
		Event: event,
	}

	request := reqbuilder.NewRequestWithBody(http.MethodPost, "/delivery/", stateTransitionDeliveryRequest)
	request.Header.Add("Token", jwt)

	return request
}
