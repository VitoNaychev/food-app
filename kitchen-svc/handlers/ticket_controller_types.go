package handlers

import (
	"time"

	"github.com/VitoNaychev/food-app/kitchen-svc/models"
)

type StateTransitionTicketRequest struct {
	ID      int                `validate:"required,min=1"    json:"id"`
	Event   models.TicketEvent `validate:"required,min=0"    json:"event"`
	ReadyBy string             `                             json:"ready_by"`
}

type StateTransitionResponse struct {
	ID    int    `validate:"required,min=1"    json:"id"`
	State string `validate:"required"          json:"state"`
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
	ID      int                     `validate:"required,min=1"    json:"id"`
	Total   float32                 `validate:"required,min=0"    json:"total"`
	State   string                  `validate:"required"          json:"state"`
	Items   []GetTicketItemResponse `validate:"required"          json:"items"`
	ReadyBy time.Time               `validate:"required"          json:"ready_by"`
}

type GetTicketItemResponse struct {
	Quantity int    `validate:"required,min=1"    json:"quantity"`
	Name     string `validate:"required"          json:"name"`
}
