package handlers

import (
	"net/http"

	"github.com/VitoNaychev/food-app/kitchen-svc/models"
	"github.com/VitoNaychev/food-app/reqbuilder"
)

func NewChangeTicketStateRequest(jwt string, id int, event models.TicketEvent) *http.Request {
	ticketRequest := StateTransitionTicketRequest{
		ID:    id,
		Event: event,
	}

	request := reqbuilder.NewRequestWithBody(http.MethodPost, "/tickets", ticketRequest)
	request.Header.Add("Token", jwt)

	return request
}

func NewGetTicketsRequest(jwt string, queryParams string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/tickets"+queryParams, nil)
	request.Header.Add("Token", jwt)

	return request
}
