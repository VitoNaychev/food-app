package models

type TicketStore interface {
	GetTicketsByRestaurantIDWhereState(int, TicketState) ([]Ticket, error)
}
