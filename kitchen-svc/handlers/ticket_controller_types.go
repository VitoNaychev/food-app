package handlers

import (
	"time"

	"github.com/VitoNaychev/food-app/kitchen-svc/models"
)

type StateTransitionTicketRequest struct {
	ID      int
	Event   models.TicketEvent
	ReadyBy string
}

type StateTransitionResponse struct {
	ID    int
	State string
}

func NewStateTransitionResponse(ticket models.Ticket) StateTransitionResponse {
	stateName, _ := models.StateValueToStateName(ticket.State)
	stateTransitionResponse := StateTransitionResponse{
		ID:    ticket.ID,
		State: stateName,
	}

	return stateTransitionResponse
}

type GetTicketResponse struct {
	ID      int
	Total   float32
	State   string
	Items   []GetTicketItemResponse
	ReadyBy time.Time
}

type GetTicketItemResponse struct {
	Quantity int
	Name     string
}
