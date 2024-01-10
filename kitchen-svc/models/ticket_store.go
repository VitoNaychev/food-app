package models

type TicketStore interface {
	GetTicketsByRestaurantIDWhereState(int, TicketState) ([]Ticket, error)
	UpdateTicketState(int, TicketState) error
	GetTicketByID(int) (Ticket, error)
}
