package handlers

import "github.com/VitoNaychev/food-app/kitchen-svc/models"

type StateTransitionTicketRequest struct {
	ID    int
	Event models.TicketEvent
}

type GetTicketResponse struct {
	ID    int
	Total float32
	Items []GetTicketItemResponse
}

type GetTicketItemResponse struct {
	Quantity int
	Name     string
}
